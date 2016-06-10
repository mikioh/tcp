// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd

package tcp_test

import (
	"strconv"
	"strings"
	"syscall"
)

func kernelVersion() []int {
	s, err := syscall.Sysctl("kern.osrelease")
	if err != nil {
		return nil
	}
	ss := strings.Split(s, ".")
	vers := make([]int, len(ss))
	for i, s := range ss {
		vers[i], _ = strconv.Atoi(s)
	}
	return vers
}
