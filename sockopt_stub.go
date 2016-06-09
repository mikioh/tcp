// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package tcp

import "time"

func setCork(s uintptr, on bool) error                           { return errOpNoSupport }
func setKeepAliveIdleInterval(s uintptr, d time.Duration) error  { return errOpNoSupport }
func setKeepAliveProbeInterval(s uintptr, d time.Duration) error { return errOpNoSupport }
func setKeepAliveProbeCount(s uintptr, probes int) error         { return errOpNoSupport }
func getInt(s uintptr, opt *sockOpt) (int, error)                { return 0, errOpNoSupport }
func setInt(s uintptr, opt *sockOpt, v int) error                { return errOpNoSupport }
