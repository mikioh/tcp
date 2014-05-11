// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"syscall"
	"unsafe"
)

func (opt *opt) setCorkedBuffer(on bool) error {
	fd, err := opt.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, ianaProtocolTCP, sysSockoptTCPNoPush, boolint(on)))
}

func (opt *opt) info() (*Info, error) {
	fd, err := opt.sysfd()
	if err != nil {
		return nil, err
	}
	var v sysTCPInfo
	l := sysSockoptLen(sysSizeofTCPInfo)
	if err := getsockopt(fd, ianaProtocolTCP, sysSockoptTCPInfo, unsafe.Pointer(&v), &l); err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return parseSysTCPInfo(&v), nil
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseSysTCPInfo(sti *sysTCPInfo) *Info {
	ti := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPInfoOptWindowScale != 0 {
		ti.Options = append(ti.Options, WindowScale(sti.Sndrcv_wscale>>4))
		ti.PeerOptions = append(ti.PeerOptions, WindowScale(sti.Sndrcv_wscale&0x0f))
	}
	return ti
}
