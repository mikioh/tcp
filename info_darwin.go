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
	SenderWindow   uint `json:"snd wnd"` // advertised sender window in bytes
	ReceiverWindow uint `json:"rcv wnd"` // advertised receiver window in bytes
}

func info(s int) (*Info, error) {
	var v sysTCPConnInfo
	l := uint32(sysSizeofTCPConnInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_CONNECTION_INFO, unsafe.Pointer(&v), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseInfo(&v), nil
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
	ti.RTT = time.Duration(sti.Rttcur) * time.Microsecond
	ti.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	ti.RTO = time.Duration(sti.Rto) * time.Microsecond
	ti.CC = &CongestionControl{
		SenderMSS:         MaxSegSize(sti.Maxseg),
		ReceiverMSS:       MaxSegSize(sti.Maxseg),
		SenderSSThreshold: uint(sti.Snd_ssthresh),
		SenderWindow:      uint(sti.Snd_cwnd),
	}
	ti.SysInfo = &SysInfo{
		SenderWindow:   uint(sti.Snd_wnd),
		ReceiverWindow: uint(sti.Rcv_wnd),
	}
	return ti
}
