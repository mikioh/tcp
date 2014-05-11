// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

const (
	// See /usr/include/netinet/tcp.h.
	sysSockoptTCPNoPush = 0x4
	sysSockoptTCPInfo   = 0x20

	sysTCPInfoOptTimestamps  = 0x1
	sysTCPInfoOptSACK        = 0x2
	sysTCPInfoOptWindowScale = 0x4
	sysTCPInfoOptECN         = 0x8
	sysTCPInfoOptTOE         = 0x10

	sysSizeofTCPInfo = 0xec
)

type sysTCPInfo struct {
	State           uint8
	_Ca_state       uint8
	_Retransmits    uint8
	_Probes         uint8
	_Backoff        uint8
	Options         uint8
	Sndrcv_wscale   uint8
	Pad_0           uint8
	Rto             uint32
	_Ato            uint32
	Snd_mss         uint32
	Rcv_mss         uint32
	_Unacked        uint32
	_Sacked         uint32
	_Lost           uint32
	_Retrans        uint32
	_Fackets        uint32
	_Last_data_sent uint32
	_Last_ack_sent  uint32
	Last_data_recv  uint32
	_Last_ack_recv  uint32
	_Pmtu           uint32
	_Rcv_ssthresh   uint32
	Rtt             uint32
	Rttvar          uint32
	Snd_ssthresh    uint32
	Snd_cwnd        uint32
	_Advmss         uint32
	_Reordering     uint32
	_Rcv_rtt        uint32
	Rcv_space       uint32
	Snd_wnd         uint32
	Snd_bwnd        uint32
	Snd_nxt         uint32
	Rcv_nxt         uint32
	Toe_tid         uint32
	Snd_rexmitpack  uint32
	Rcv_ooopack     uint32
	Snd_zerowin     uint32
	Pad             [26]uint32
}
