// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!freebsd,!linux,!netbsd,!windows

package tcp

// A SysInfo represents platform-specific information.
type SysInfo struct{}

func info(s int) (*Info, error) { return nil, errOpNoSupport }

// A SysCongestionControl represents platform-specific congestion
// control information.
type SysCongestionControl struct{}
