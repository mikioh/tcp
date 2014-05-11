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
	AdvertisedMSS  MaxSegSize // advertised maximum segment size
	ReceiverWindow uint       // advertised receiver window in bytes
	CAState        CAState    // state of congestion avoidance
	UnackSegs      uint       // # of unack'd segments in transmission queue
	SackSegs       uint       // # of sack'd segments in tranmission queue
	LostSegs       uint       // # of lost segments in transmission queue
	RetransSegs    uint       // # of retransmitted segments in transmission queue
	ForwardAckSegs uint       // # of forward ack'd segments in transmission queue
}
