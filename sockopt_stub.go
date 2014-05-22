// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package tcp

func (opt *opt) setMaxKeepAliveProbes(max int) error {
	return errOpNoSupport
}

func (opt *opt) setCork(on bool) error {
	return errOpNoSupport
}

func (opt *opt) info() (*Info, error) {
	return nil, errOpNoSupport
}
