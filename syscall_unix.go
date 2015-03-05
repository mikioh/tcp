// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux,amd64 linux,arm linux,ppc64 linux,ppc64le netbsd openbsd

package tcp

import (
	"syscall"
	"unsafe"
)

func getsockopt(fd int, level, name int, v unsafe.Pointer, l *sysSockoptLen) error {
	if _, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(fd), uintptr(level), uintptr(name), uintptr(v), uintptr(unsafe.Pointer(l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}

func getsockoptIntByIoctl(fd, ioc int) (int, error) {
	var i int
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(ioc), uintptr(unsafe.Pointer(&i))); errno != 0 {
		return 0, error(errno)
	}
	return i, nil
}
