// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package tcp

import (
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

func buffered(s uintptr) int {
	o, ok := soOptions[soBuffered]
	if !ok {
		return -1
	}
	var b [4]byte
	if err := ioctl(s, o.name, b[:]); err != nil {
		return -1
	}
	return int(nativeEndian.Uint32(b[:]))
}

func available(s uintptr) int {
	o, ok := soOptions[soAvailable]
	if !ok {
		return -1
	}
	var b [4]byte
	if err := ioctl(s, o.name, b[:]); err != nil {
		return -1
	}
	n := int(nativeEndian.Uint32(b[:]))
	if runtime.GOOS == "linux" {
		var o tcpopt.SendBuffer
		l := uint32(4)
		if err := getsockopt(s, o.Level(), o.Name(), b[:], &l); err != nil {
			return -1
		}
		return int(nativeEndian.Uint32(b[:])) - n
	}
	return n
}

func setCork(s uintptr, on bool) error {
	o, ok := soOptions[soCork]
	if !ok {
		return errOpNoSupport
	}
	v := boolint32(bool(on))
	b := (*[4]byte)(unsafe.Pointer(&v))[:]
	return os.NewSyscallError("setsockopt", setsockopt(s, o.level, o.name, b, 4))
}

func setUnsentThreshold(s uintptr, n int) error {
	o, ok := soOptions[soUnsentThreshold]
	if !ok {
		return errOpNoSupport
	}
	v := int32(n)
	b := (*[4]byte)(unsafe.Pointer(&v))[:]
	return os.NewSyscallError("setsockopt", setsockopt(s, o.level, o.name, b, 4))
}

func setKeepAliveIdleInterval(s uintptr, d time.Duration) error {
	o := tcpopt.KeepAliveIdleInterval(d)
	b, err := o.Marshal()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", setsockopt(s, o.Level(), o.Name(), b, uint32(len(b))))
}

func setKeepAliveProbeInterval(s uintptr, d time.Duration) error {
	o := tcpopt.KeepAliveProbeInterval(d)
	b, err := o.Marshal()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", setsockopt(s, o.Level(), o.Name(), b, uint32(len(b))))
}

func setKeepAliveProbeCount(s uintptr, n int) error {
	o := tcpopt.KeepAliveProbeCount(n)
	b, err := o.Marshal()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", setsockopt(s, o.Level(), o.Name(), b, uint32(len(b))))
}
