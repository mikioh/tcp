// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !freebsd,!linux

package tcp

func (opt *opt) setCorkedBuffer(on bool) error {
	return errOpNoSupport
}

func (opt *opt) info() (*Info, error) {
	return nil, errOpNoSupport
}
