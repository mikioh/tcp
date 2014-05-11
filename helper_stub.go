// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package tcp

func (c *opt) sysfd() (int, error) {
	return 0, errOpNoSupport
}
