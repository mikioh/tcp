// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package tcp

import "time"

func setCork(s int, on bool) error                           { return errOpNoSupport }
func setKeepAliveIdleInterval(s int, d time.Duration) error  { return errOpNoSupport }
func setKeepAliveProbeInterval(s int, d time.Duration) error { return errOpNoSupport }
func setKeepAliveProbeCount(s int, probes int) error         { return errOpNoSupport }
func getInt(s int, opt *sockOpt) (int, error)                { return 0, errOpNoSupport }
func setInt(s int, opt *sockOpt, v int) error                { return errOpNoSupport }
