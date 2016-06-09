// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd netbsd

package tcp

import (
	"os"
	"time"
	"unsafe"
)

// A SysInfo represents platform-specific information.
type SysInfo struct {
	SenderWindow      uint `json:"snd wnd"`         // advertised sender window in bytes
	NextEgressSeq     uint `json:"egress seq"`      // next egress seq. number
	NextIngressSeq    uint `json:"ingress seq"`     // next ingress seq. number
	RetransSegs       uint `json:"retrans segs"`    // # of retransmit segments sent
	OutOfOrderSegs    uint `json:"ooo segs"`        // # of out-of-order segments received
	ZeroWindowUpdates uint `json:"zerownd updates"` // # of zero-window updates sent
	Offloading        bool `json:"offloading"`      // TCP offload processing
}

func info(s uintptr) (*Info, error) {
	var sti sysTCPInfo
	l := uint32(sizeofTCPInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_INFO, unsafe.Pointer(&sti), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseInfo(&sti), nil
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseInfo(sti *sysTCPInfo) *Info {
	ti := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPI_OPT_WSCALE != 0 {
		ti.Options = append(ti.Options, WindowScale(sti.Pad_cgo_0[0]>>4))
		ti.PeerOptions = append(ti.PeerOptions, WindowScale(sti.Pad_cgo_0[0]&0x0f))
	}
	if sti.Options&sysTCPI_OPT_TIMESTAMPS != 0 {
		ti.Options = append(ti.Options, Timestamps(true))
		ti.PeerOptions = append(ti.PeerOptions, Timestamps(true))
	}
	ti.SenderMSS = MaxSegSize(sti.Snd_mss)
	ti.ReceiverMSS = MaxSegSize(sti.Rcv_mss)
	ti.RTT = time.Duration(sti.Rtt) * time.Microsecond
	ti.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	ti.RTO = time.Duration(sti.Rto) * time.Microsecond
	ti.ATO = time.Duration(sti.X__tcpi_ato) * time.Microsecond
	ti.LastDataSent = time.Duration(sti.X__tcpi_last_data_sent) * time.Microsecond
	ti.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Microsecond
	ti.LastAckReceived = time.Duration(sti.X__tcpi_last_ack_recv) * time.Microsecond
	ti.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_space),
	}
	ti.CongestionControl = &CongestionControl{
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.X__tcpi_rcv_ssthresh),
		SenderWindow:        uint(sti.Snd_cwnd),
	}
	ti.Sys = &SysInfo{
		SenderWindow:      uint(sti.Snd_wnd),
		NextEgressSeq:     uint(sti.Snd_nxt),
		NextIngressSeq:    uint(sti.Rcv_nxt),
		RetransSegs:       uint(sti.Snd_rexmitpack),
		OutOfOrderSegs:    uint(sti.Rcv_ooopack),
		ZeroWindowUpdates: uint(sti.Snd_zerowin),
	}
	if sti.Options&sysTCPI_OPT_TOE != 0 {
		ti.Sys.Offloading = true
	}
	return ti
}

// A SysCongestionControl represents platform-specific congestion
// control information.
type SysCongestionControl struct{}
