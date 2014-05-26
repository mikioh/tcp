// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"encoding/json"
	"time"
)

var _ json.Marshaler = &Info{}

// A State represents a state of TCP connection.
type State int

const (
	Unknown State = iota
	Closed
	Listen
	SynSent
	SynReceived
	Established
	FinWait1
	FinWait2
	CloseWait
	LastAck
	Closing
	TimeWait
)

var states = map[State]string{
	Unknown:     "unknown",
	Closed:      "closed",
	Listen:      "listen",
	SynSent:     "syn-sent",
	SynReceived: "syn-received",
	Established: "established",
	FinWait1:    "fin-wait-1",
	FinWait2:    "fin-wait-2",
	CloseWait:   "close-wait",
	LastAck:     "last-ack",
	Closing:     "closing",
	TimeWait:    "time-wait",
}

func (st State) String() string {
	s, ok := states[st]
	if !ok {
		return "<nil>"
	}
	return s
}

// An Info represents TCP connection information.
type Info struct {
	State            State              `json:"state"`               // connection state
	Options          []Option           `json:"opts,omitempty"`      // requesting options
	PeerOptions      []Option           `json:"peer opts,omitempty"` // options requested from peer
	RTT              time.Duration      `json:"rtt"`                 // round-trip time
	RTTVar           time.Duration      `json:"rtt var"`             // round-trip time variation
	RTO              time.Duration      `json:"rto"`                 // retransmission timeout
	ATO              time.Duration      `json:"ato"`                 // delayed acknowledgement timeout [linux only]
	LastDataSent     time.Duration      `json:"last data sent"`      // since last data sent [linux only]
	LastDataReceived time.Duration      `json:"last data rcvd"`      // since last data received
	LastAckReceived  time.Duration      `json:"last ack rcvd"`       // since last ack received [linux only]
	CC               *CongestionControl `json:"cc,omitempty"`        // congestion control information
	SysInfo          *SysInfo           `json:"sys info,omitempty"`  // platform-specific information
}

// MarshalJSON implements the MarshalJSON method of json.Marshaler
// interface.
func (info *Info) MarshalJSON() ([]byte, error) {
	raw := make(map[string]interface{})
	raw["state"] = info.State.String()
	if len(info.Options) > 0 {
		opts := make(map[string]interface{})
		for _, opt := range info.Options {
			opts[opt.Kind().String()] = opt
		}
		raw["opts"] = opts
	}
	if len(info.PeerOptions) > 0 {
		opts := make(map[string]interface{})
		for _, opt := range info.PeerOptions {
			opts[opt.Kind().String()] = opt
		}
		raw["peer opts"] = opts
	}
	raw["rtt"] = info.RTT
	raw["rtt var"] = info.RTTVar
	raw["rto"] = info.RTO
	raw["ato"] = info.ATO
	raw["last data sent"] = info.LastDataSent
	raw["last data rcvd"] = info.LastDataReceived
	raw["last ack rcvd"] = info.LastAckReceived
	if info.CC != nil {
		raw["cc"] = info.CC
	}
	if info.SysInfo != nil {
		raw["sys info"] = info.SysInfo
	}
	return json.Marshal(&raw)
}

// A CongestionControl represents TCP congestion control information.
type CongestionControl struct {
	SenderMSS           MaxSegSize `json:"snd mss"`      // maximum segment size for sender
	ReceiverMSS         MaxSegSize `json:"rcv mss"`      // maximum sengment size for receiver
	SenderSSThreshold   uint       `json:"snd ssthresh"` // slow start threshold for sender
	ReceiverSSThreshold uint       `json:"rcv ssthresh"` // slow start threshold for receiver [linux only]
	SenderWindow        uint       `json:"cwnd"`         // congestion window for sender
}
