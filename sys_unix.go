// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux,amd64 linux,arm linux,armbe linux,arm64 linux,ppc64 linux,ppc64le linux,mips linux,mipsle linux,mips64 linux,mips64le netbsd openbsd

package tcp

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

func buffered(s uintptr) int {
	var b [4]byte
	if err := ioctl(s, options[soBuffered].name, b[:]); err != nil {
		return -1
	}
	return int(nativeEndian.Uint32(b[:]))
}

func available(s uintptr) int {
	var b [4]byte
	if runtime.GOOS == "darwin" {
		if err := getsockopt(s, options[soAvailable].level, options[soAvailable].name, b[:]); err != nil {
			return -1
		}
	} else {
		if err := ioctl(s, options[soAvailable].name, b[:]); err != nil {
			return -1
		}
	}
	n := int(nativeEndian.Uint32(b[:]))
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		var o tcpopt.SendBuffer
		if err := getsockopt(s, o.Level(), o.Name(), b[:]); err != nil {
			return -1
		}
		return int(nativeEndian.Uint32(b[:])) - n
	}
	return n
}

func ioctl(s uintptr, ioc int, b []byte) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s, uintptr(ioc), uintptr(unsafe.Pointer(&b[0]))); errno != 0 {
		return error(errno)
	}
	return nil
}

func setsockopt(s uintptr, level, name int, b []byte) error {
	l := uint32(len(b))
	if _, _, errno := syscall.Syscall6(syscall.SYS_SETSOCKOPT, s, uintptr(level), uintptr(name), uintptr(unsafe.Pointer(&b[0])), uintptr(l), 0); errno != 0 {
		return error(errno)
	}
	return nil
}

func getsockopt(s uintptr, level, name int, b []byte) error {
	l := uint32(len(b))
	if _, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, s, uintptr(level), uintptr(name), uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}
