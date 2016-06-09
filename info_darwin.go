// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"time"
	"unsafe"
)

// A SysInfo represents platform-specific information.
type SysInfo struct {
	Flags        uint          `json:"flags"`     // flags
	SenderWindow uint          `json:"snd wnd"`   // advertised sender window in bytes
	SenderInUse  uint          `json:"snd inuse"` // bytes in send buffer including inflight data
	SRTT         time.Duration `json:"srtt"`      // smoothed round-trip time
}

func info(s uintptr) (*Info, error) {
	var sti sysTCPConnInfo
	l := uint32(sizeofTCPConnInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_CONNECTION_INFO, unsafe.Pointer(&sti), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseInfo(&sti), nil
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseInfo(sti *sysTCPConnInfo) *Info {
	ti := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPCI_OPT_WSCALE != 0 {
		ti.Options = append(ti.Options, WindowScale(sti.Snd_wscale))
		ti.PeerOptions = append(ti.PeerOptions, WindowScale(sti.Rcv_wscale))
	}
	if sti.Options&sysTCPCI_OPT_TIMESTAMPS != 0 {
		ti.Options = append(ti.Options, Timestamps(true))
		ti.PeerOptions = append(ti.PeerOptions, Timestamps(true))
	}
	ti.SenderMSS = MaxSegSize(sti.Maxseg)
	ti.ReceiverMSS = MaxSegSize(sti.Maxseg)
	ti.RTT = time.Duration(sti.Rttcur) * time.Millisecond
	ti.RTTVar = time.Duration(sti.Rttvar) * time.Millisecond
	ti.RTO = time.Duration(sti.Rto) * time.Millisecond
	ti.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_wnd),
	}
	ti.CongestionControl = &CongestionControl{
		SenderSSThreshold: uint(sti.Snd_ssthresh),
		SenderWindow:      uint(sti.Snd_cwnd),
	}
	ti.Sys = &SysInfo{
		Flags:        uint(sti.Flags),
		SenderWindow: uint(sti.Snd_wnd),
		SenderInUse:  uint(sti.Snd_sbbytes),
		SRTT:         time.Duration(sti.Srtt) * time.Millisecond,
	}
	return ti
}

// A SysCongestionControl represents platform-specific congestion
// control information.
type SysCongestionControl struct{}
