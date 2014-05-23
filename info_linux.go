// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

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
	AdvertisedMSS   MaxSegSize `json:"adv mss"'`     // advertised maximum segment size
	ReceiverWindow  uint       `json:"rcv wnd"`      // advertised receiver window in bytes
	CAState         CAState    `json:"ca state"`     // state of congestion avoidance
	KeepAliveProbes uint       `json:"ka probes"`    // # of keep alive probes sent"`
	UnackSegs       uint       `json:"unack segs"`   // # of unack'd segments in transmission queue
	SackSegs        uint       `json:"sack segs"`    // # of sack'd segments in tranmission queue
	LostSegs        uint       `json:"lost segs"`    // # of lost segments in transmission queue
	RetransSegs     uint       `json:"retrans segs"` // # of retransmitting segments in transmission queue
	ForwardAckSegs  uint       `json:"fack segs"`    // # of forward ack'd segments in transmission queue
}
