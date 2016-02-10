// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "syscall"

func buffered(s syscall.Handle) int {
	return -1
}

func available(s syscall.Handle) int {
	return -1
}
