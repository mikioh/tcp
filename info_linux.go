// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"time"
	"unsafe"
)

// A CAState represents a state of congestion avoidance.
type CAState int

var caStates = map[CAState]string{
	CAOpen:     "open",
	CADisorder: "disorder",
	CACWR:      "congestion window reduced",
	CARecovery: "recovery",
	CALoss:     "loss",
}

func (st CAState) String() string {
	s, ok := caStates[st]
	if !ok {
		return "<nil>"
	}
	return s
}

// A SysInfo represents platform-specific information.
type SysInfo struct {
	PathMTU         uint       `json:"path mtu"`     // path maximum transmission unit
	AdvertisedMSS   MaxSegSize `json:"adv mss"`      // advertised maximum segment size
	ReceiverWindow  uint       `json:"rcv wnd"`      // advertised receiver window in bytes
	CAState         CAState    `json:"ca state"`     // state of congestion avoidance
	KeepAliveProbes uint       `json:"ka probes"`    // # of keep alive probes sent
	UnackSegs       uint       `json:"unack segs"`   // # of unack'd segments in transmission queue
	SackSegs        uint       `json:"sack segs"`    // # of sack'd segments in tranmission queue
	LostSegs        uint       `json:"lost segs"`    // # of lost segments in transmission queue
	RetransSegs     uint       `json:"retrans segs"` // # of retransmitting segments in transmission queue
	ForwardAckSegs  uint       `json:"fack segs"`    // # of forward ack'd segments in transmission queue
}

func info(s int) (*Info, error) {
	var v sysTCPInfo
	l := uint32(sysSizeofTCPInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_INFO, unsafe.Pointer(&v), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseInfo(&v), nil
}

var sysStates = [12]State{Unknown, Established, SynSent, SynReceived, FinWait1, FinWait2, TimeWait, Closed, CloseWait, LastAck, Listen, Closing}

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
	ti.RTT = time.Duration(sti.Rtt) * time.Microsecond
	ti.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	ti.RTO = time.Duration(sti.Rto) * time.Microsecond
	ti.ATO = time.Duration(sti.Ato) * time.Microsecond
	ti.LastDataSent = time.Duration(sti.Last_data_sent) * time.Millisecond
	ti.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Millisecond
	ti.LastAckReceived = time.Duration(sti.Last_ack_recv) * time.Millisecond
	ti.CC = &CongestionControl{
		SenderMSS:           MaxSegSize(sti.Snd_mss),
		ReceiverMSS:         MaxSegSize(sti.Rcv_mss),
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.Rcv_ssthresh),
		SenderWindow:        uint(sti.Snd_cwnd),
	}
	ti.SysInfo = &SysInfo{
		PathMTU:         uint(sti.Pmtu),
		AdvertisedMSS:   MaxSegSize(sti.Advmss),
		ReceiverWindow:  uint(sti.Rcv_space),
		CAState:         CAState(sti.Ca_state),
		KeepAliveProbes: uint(sti.Probes),
		UnackSegs:       uint(sti.Unacked),
		SackSegs:        uint(sti.Sacked),
		LostSegs:        uint(sti.Lost),
		RetransSegs:     uint(sti.Retrans),
		ForwardAckSegs:  uint(sti.Fackets),
	}
	return ti
}
