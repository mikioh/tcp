// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.4

package tcp

import (
	"syscall"
	"unsafe"
)

func socketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (int, syscall.Errno)

func getsockopt(s, level, name int, v unsafe.Pointer, l *sysSockoptLen) error {
	if _, errno := socketcall(sysGETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(v), uintptr(unsafe.Pointer(l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}
