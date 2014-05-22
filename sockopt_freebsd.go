// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

func (opt *opt) setMaxKeepAliveProbes(max int) error {
	fd, err := opt.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, max))
}

func (opt *opt) setCork(on bool) error {
	fd, err := opt.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NOPUSH, boolint(on)))
}

func (opt *opt) info() (*Info, error) {
	fd, err := opt.sysfd()
	if err != nil {
		return nil, err
	}
	var v sysTCPInfo
	l := sysSockoptLen(sysSizeofTCPInfo)
	if err := getsockopt(fd, syscall.IPPROTO_TCP, sysSockoptTCPInfo, unsafe.Pointer(&v), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseTCPInfo(&v), nil
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseTCPInfo(sti *sysTCPInfo) *Info {
	ti := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPIOptWscale != 0 {
		ti.Options = append(ti.Options, WindowScale(sti.Pad_cgo_0[0]>>4))
		ti.PeerOptions = append(ti.PeerOptions, WindowScale(sti.Pad_cgo_0[0]&0x0f))
	}
	if sti.Options&sysTCPIOptTimestamps != 0 {
		ti.Options = append(ti.Options, Timestamps(true))
		ti.PeerOptions = append(ti.PeerOptions, Timestamps(true))
	}
	ti.SenderMSS = MaxSegSize(sti.Snd_mss)
	ti.ReceiverMSS = MaxSegSize(sti.Rcv_mss)
	ti.LastDataSent = time.Duration(sti.X__tcpi_last_data_sent) * time.Microsecond
	ti.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Microsecond
	ti.LastAckReceived = time.Duration(sti.X__tcpi_last_ack_recv) * time.Microsecond
	ti.CC = &CongestionControl{
		RTO:                 time.Duration(sti.Rto) * time.Microsecond,
		ATO:                 time.Duration(sti.X__tcpi_ato) * time.Microsecond,
		RTT:                 time.Duration(sti.Rtt) * time.Microsecond,
		RTTStdDev:           time.Duration(sti.Rttvar) * time.Microsecond,
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.X__tcpi_rcv_ssthresh),
		SenderWindow:        uint(sti.Snd_cwnd),
	}
	ti.SysInfo = &SysInfo{
		SenderWindow:      uint(sti.Snd_wnd),
		ReceiverWindow:    uint(sti.Rcv_space),
		NextEgressSeq:     uint(sti.Snd_nxt),
		NextIngressSeq:    uint(sti.Rcv_nxt),
		RetransSegs:       uint(sti.Snd_rexmitpack),
		OutOfOrderSegs:    uint(sti.Rcv_ooopack),
		ZeroWindowUpdates: uint(sti.Snd_zerowin),
	}
	if sti.Options&sysTCPIOptTOE != 0 {
		ti.SysInfo.Offloading = true
	}
	return ti
}
