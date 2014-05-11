// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func (opt *opt) info() (*Info, error) {
	return nil, errOpNoSupport
}
