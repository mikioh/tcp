// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package tcp

import "errors"

var options [soMax]option

func buffered(s uintptr) int  { return -1 }
func available(s uintptr) int { return -1 }

func setsockopt(s uintptr, level, name int, b []byte) error {
	return errors.New("operation not supported")
}

func getsockopt(s uintptr, level, name int, b []byte) error {
	return errors.New("operation not supported")
}
