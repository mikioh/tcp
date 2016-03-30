// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux,amd64 linux,arm linux,armbe linux,arm64 linux,ppc64 linux,ppc64le linux,mips linux,mipsle linux,mips64 linux,mips64le netbsd openbsd

package tcp

import (
	"syscall"
	"unsafe"
)

func getsockopt(s int, level, name int, v unsafe.Pointer, l *uint32) error {
	if _, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(v), uintptr(unsafe.Pointer(l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}

func ioctl(s, ioc int) (int, error) {
	var i int
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s), uintptr(ioc), uintptr(unsafe.Pointer(&i))); errno != 0 {
		return 0, error(errno)
	}
	return i, nil
}
