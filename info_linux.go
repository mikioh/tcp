// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"strings"
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
	CAState         CAState    `json:"ca state"`     // state of congestion avoidance
	KeepAliveProbes uint       `json:"ka probes"`    // # of keep alive probes sent
	UnackSegs       uint       `json:"unack segs"`   // # of unack'd segments in transmission queue
	SackSegs        uint       `json:"sack segs"`    // # of sack'd segments in tranmission queue
	LostSegs        uint       `json:"lost segs"`    // # of lost segments in transmission queue
	RetransSegs     uint       `json:"retrans segs"` // # of retransmitting segments in transmission queue
	ForwardAckSegs  uint       `json:"fack segs"`    // # of forward ack'd segments in transmission queue
}

func info(s int) (*Info, error) {
	var sti sysTCPInfo
	l := uint32(sysSizeofTCPInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_INFO, unsafe.Pointer(&sti), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	var stcci sysTCPCCInfo
	l = uint32(sysSizeofTCPCCInfo)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_CC_INFO, unsafe.Pointer(&stcci), &l); err != nil {
		return parseInfo(&sti), nil
	}
	b := make([]byte, 16) // see TCP_CA_NAME_MAX
	l = uint32(16)
	if err := getsockopt(s, ianaProtocolTCP, sysTCP_CONGESTION, unsafe.Pointer(&b[0]), &l); err != nil {
		return parseInfo(&sti), nil
	}
	ti := parseInfo(&sti)
	i := 0
	for i = 0; i < 16; i++ {
		if b[i] == 0 {
			break
		}
	}
	scc := parseSysCC(string(b[:i]), &stcci)
	if ti != nil && ti.CongestionControl != nil && scc != nil {
		ti.CongestionControl.Sys = scc
	}
	return ti, nil
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
	ti.SenderMSS = MaxSegSize(sti.Snd_mss)
	ti.ReceiverMSS = MaxSegSize(sti.Rcv_mss)
	ti.RTT = time.Duration(sti.Rtt) * time.Microsecond
	ti.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	ti.RTO = time.Duration(sti.Rto) * time.Microsecond
	ti.ATO = time.Duration(sti.Ato) * time.Microsecond
	ti.LastDataSent = time.Duration(sti.Last_data_sent) * time.Millisecond
	ti.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Millisecond
	ti.LastAckReceived = time.Duration(sti.Last_ack_recv) * time.Millisecond
	ti.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_space),
	}
	ti.CongestionControl = &CongestionControl{
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.Rcv_ssthresh),
		SenderWindow:        uint(sti.Snd_cwnd),
	}
	ti.Sys = &SysInfo{
		PathMTU:         uint(sti.Pmtu),
		AdvertisedMSS:   MaxSegSize(sti.Advmss),
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

// A SysCongestionControl represents platform-specific congestion
// control information.
type SysCongestionControl struct {
	Algo string                   `json:"algo"` // congestion control algorithm name
	Info SysCongestionControlInfo `json:"info,omitempty"`
}

// A SysCongestionControlInfo represents congestion control
// algorithm-specific information.
type SysCongestionControlInfo interface {
	String() string
}

// A VegasInfo represents TCP Vegas congestion control information.
type VegasInfo struct {
	Enabled    bool          `json:"enabled"`
	RoundTrips uint          `json:"rnd trips"` // # of round-trips
	RTT        time.Duration `json:"rtt"`       // round-trip time
	MinRTT     time.Duration `json:"min rtt"`   // minimum round-trip time
}

// String implements the String method of SysCongestionControlInfo
// interface.
func (vi *VegasInfo) String() string { return "vegas" }

// A CEState represents a state of ECN congestion encountered (CE)
// codepoint.
type CEState int

// A DCTCPInfo represents Datacenter TCP congestion control
// information.
type DCTCPInfo struct {
	Enabled         bool    `json:"enabled"`
	CEState         CEState `json:"ce state"`    // state of ECN CE codepoint
	Alpha           uint    `json:"alpha"`       // fraction of bytes sent
	ECNAckedBytes   uint    `json:"ecn acked"`   // # of acked bytes with ECN
	TotalAckedBytes uint    `json:"total acked"` // total # of acked bytes
}

// String implements the String method of SysCongestionControlInfo
// interface.
func (di *DCTCPInfo) String() string { return "dctcp" }

func parseSysCC(name string, stcci *sysTCPCCInfo) *SysCongestionControl {
	scc := SysCongestionControl{Algo: name}
	if strings.HasPrefix(name, "dctcp") {
		stdi := (*sysTCPDCTCPInfo)(unsafe.Pointer(stcci))
		if stdi.Enabled == 0 {
			return &scc
		}
		scc.Info = &DCTCPInfo{
			Enabled: true,
			Alpha:   uint(stdi.Alpha),
		}
		return &scc
	}
	stvi := (*sysTCPVegasInfo)(unsafe.Pointer(stcci))
	if stvi.Enabled == 0 {
		return &scc
	}
	scc.Info = &VegasInfo{
		Enabled:    true,
		RoundTrips: uint(stvi.Rttcnt),
		RTT:        time.Duration(stvi.Rtt) * time.Microsecond,
		MinRTT:     time.Duration(stvi.Minrtt) * time.Microsecond,
	}
	return &scc
}
