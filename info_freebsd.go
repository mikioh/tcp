// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

// A SysInfo represents platform-specific information.
type SysInfo struct {
	SenderWindow      uint // advertised sender window in bytes
	ReceiverWindow    uint // advertised receiver window in bytes
	NextEgressSeq     uint // next egress seq. number
	NextIngressSeq    uint // next ingress seq. number
	RetransSegs       uint // # of retransmit segments sent
	OutOfOrderSegs    uint // #of out-of-order segments received
	ZeroWindowUpdates uint // # of zero-window updates sent
	Offloading        bool // TCP offload processing
}
