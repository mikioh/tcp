// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

const (
	// See /usr/include/netinet/tcp.h.
	sysSockoptTCPCork = 0x3
	sysSockoptTCPInfo = 0xb

	sysTCPInfoOptTimestamps  = 0x1
	sysTCPInfoOptSACK        = 0x2
	sysTCPInfoOptWindowScale = 0x4
	sysTCPInfoOptECN         = 0x8
	sysTCPInfoOptECNSeen     = 0x10
	sysTCPInfoOptSynData     = 0x20

	sysSizeofTCPInfo = 0x68
)

type sysTCPInfo struct {
	State          uint8
	Ca_state       uint8
	Retransmits    uint8
	Probes         uint8
	Backoff        uint8
	Options        uint8
	Sndrcv_wscale  uint8
	Pad_0          uint8
	Rto            uint32
	Ato            uint32
	Snd_mss        uint32
	Rcv_mss        uint32
	Unacked        uint32
	Sacked         uint32
	Lost           uint32
	Retrans        uint32
	Fackets        uint32
	Last_data_sent uint32
	Last_ack_sent  uint32
	Last_data_recv uint32
	Last_ack_recv  uint32
	Pmtu           uint32
	Rcv_ssthresh   uint32
	Rtt            uint32
	Rttvar         uint32
	Snd_ssthresh   uint32
	Snd_cwnd       uint32
	Advmss         uint32
	Reordering     uint32
	Rcv_rtt        uint32
	Rcv_space      uint32
	Total_retrans  uint32
}
