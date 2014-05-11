// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

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

// An Info represents a TCP information.
type Info struct {
	State       State    // connection state
	Options     []Option // requesting options
	PeerOptions []Option // requested options
}
