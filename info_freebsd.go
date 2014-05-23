// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

// A SysInfo represents platform-specific information.
type SysInfo struct {
	SenderWindow      uint `json:"snd wnd"`         // advertised sender window in bytes
	ReceiverWindow    uint `json:"rcv wnd"`         // advertised receiver window in bytes
	NextEgressSeq     uint `json:"egress seq"`      // next egress seq. number
	NextIngressSeq    uint `json:"ingress seq"`     // next ingress seq. number
	RetransSegs       uint `json:"retrans segs"`    // # of retransmit segments sent
	OutOfOrderSegs    uint `json:"ooo segs"`        // # of out-of-order segments received
	ZeroWindowUpdates uint `json:"zerownd updates"` // # of zero-window updates sent
	Offloading        bool `json:"offloading"`      // TCP offload processing
}
