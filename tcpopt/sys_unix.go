// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package tcpopt

import (
	"time"
	"unsafe"
)

// Marshal implements the Marshal method of Option interface.
func (nd NoDelay) Marshal() ([]byte, error) {
	v := boolint32(bool(nd))
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (sb SendBuffer) Marshal() ([]byte, error) {
	v := int32(sb)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (rb ReceiveBuffer) Marshal() ([]byte, error) {
	v := int32(rb)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (ka KeepAlive) Marshal() ([]byte, error) {
	v := boolint32(bool(ka))
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (ka KeepAliveIdleInterval) Marshal() ([]byte, error) {
	ka += KeepAliveIdleInterval(options[kaIdleInterval].uot - time.Nanosecond)
	v := int32(time.Duration(ka) / options[kaIdleInterval].uot)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (ka KeepAliveProbeInterval) Marshal() ([]byte, error) {
	ka += KeepAliveProbeInterval(options[kaProbeInterval].uot - time.Nanosecond)
	v := int32(time.Duration(ka) / options[kaProbeInterval].uot)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (ka KeepAliveProbeCount) Marshal() ([]byte, error) {
	v := int32(ka)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

// Marshal implements the Marshal method of Option interface.
func (bt BufferUnsentThreshold) Marshal() ([]byte, error) {
	v := int32(bt)
	return (*[4]byte)(unsafe.Pointer(&v))[:], nil
}

func parseNoDelay(b []byte) (Option, error) {
	return NoDelay(uint32bool(nativeEndian.Uint32(b))), nil
}

func parseSendBuffer(b []byte) (Option, error) {
	return SendBuffer(nativeEndian.Uint32(b)), nil
}

func parseReceiveBuffer(b []byte) (Option, error) {
	return ReceiveBuffer(nativeEndian.Uint32(b)), nil
}

func parseKeepAlive(b []byte) (Option, error) {
	return KeepAlive(uint32bool(nativeEndian.Uint32(b))), nil
}

func parseKeepAliveIdleInterval(b []byte) (Option, error) {
	v := time.Duration(nativeEndian.Uint32(b)) * options[kaIdleInterval].uot
	return KeepAliveIdleInterval(v), nil
}

func parseKeepAliveProbeInterval(b []byte) (Option, error) {
	v := time.Duration(nativeEndian.Uint32(b)) * options[kaProbeInterval].uot
	return KeepAliveProbeInterval(v), nil
}

func parseKeepAliveProbeCount(b []byte) (Option, error) {
	return KeepAliveProbeCount(nativeEndian.Uint32(b)), nil
}

func parseBufferUnsentThreshold(b []byte) (Option, error) {
	return BufferUnsentThreshold(nativeEndian.Uint32(b)), nil
}
