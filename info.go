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
	State            State              `json:"state"`          // connection state
	Options          []Option           `json:"opts"`           // requesting options
	PeerOptions      []Option           `json:"peer opts"`      // options requested from peer
	SenderMSS        MaxSegSize         `json:"snd mss"`        // maximum segment size for sender
	ReceiverMSS      MaxSegSize         `json:"rcv mss"`        // maximum sengment size for receiver
	LastDataSent     time.Duration      `json:"last data sent"` // since last data sent [linux only]
	LastDataReceived time.Duration      `json:"last data rcvd"` // since last data received
	LastAckReceived  time.Duration      `json:"last ack rcvd"`  // since last ack received [linux only]
	CC               *CongestionControl `json:"cc"`             // congestion control information
	SysInfo          *SysInfo           `json:"sysinfo"`        // platform-specific information
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
	raw["snd mss"] = info.SenderMSS
	raw["rcv mss"] = info.ReceiverMSS
	raw["last data sent"] = info.LastDataSent
	raw["last data rcvd"] = info.LastDataReceived
	raw["last ack rcvd"] = info.LastAckReceived
	if info.CC != nil {
		raw["cc"] = info.CC
	}
	if info.SysInfo != nil {
		raw["sysinfo"] = info.SysInfo
	}
	return json.Marshal(&raw)
}

// A CongestionControl represents TCP congestion control information.
type CongestionControl struct {
	RTO                 time.Duration `json:"rto"`          // retransmission timeout
	ATO                 time.Duration `json:"ato"`          // delayed acknowledgement timeout [linux only]
	RTT                 time.Duration `json:"rtt"`          // round trip time
	RTTStdDev           time.Duration `json:"rtt stddev"`   // standard deviation of round trip time
	SenderSSThreshold   uint          `json:"snd ssthresh"` // slow start threshold for sender
	ReceiverSSThreshold uint          `json:"rcv ssthresh"` // slow start threshold for receiver [linux only]
	SenderWindow        uint          `json:"cwnd"`         // congestion window for sender
}
