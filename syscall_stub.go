// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package tcp

func setsockopt(s uintptr, level, name int, b []byte) error { return errOpNoSupport }
func getsockopt(s uintptr, level, name int, b []byte) error { return errOpNoSupport }
