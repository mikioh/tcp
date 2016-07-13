// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"encoding/binary"
	"unsafe"
)

var nativeEndian binary.ByteOrder

func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		nativeEndian = binary.LittleEndian
	} else {
		nativeEndian = binary.BigEndian
	}
}

const (
	ianaProtocolTCP = 0x6
)

const (
	soBuffered = iota
	soAvailable
	soMax
)

type option struct {
	level int // option level
	name  int // option name, must be equal or greater than 1
}
