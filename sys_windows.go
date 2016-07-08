// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var soOptions = map[int]soOption{}

func buffered(s uintptr) int                    { return -1 }
func available(s uintptr) int                   { return -1 }
func setCork(s uintptr, on bool) error          { return errOpNoSupport }
func setUnsentThreshold(s uintptr, n int) error { return errOpNoSupport }

var keepAlive = struct {
	sync.RWMutex
	syscall.TCPKeepalive
}{
	TCPKeepalive: syscall.TCPKeepalive{
		OnOff:    1,
		Time:     uint32(2 * time.Hour / time.Millisecond),
		Interval: uint32(time.Second / time.Millisecond),
	},
}

func setKeepAliveIdleInterval(s uintptr, d time.Duration) error {
	keepAlive.Lock()
	defer keepAlive.Unlock()
	d += (time.Millisecond - time.Nanosecond)
	msecs := uint32(d / time.Millisecond)
	prev := keepAlive.Time
	keepAlive.Time = msecs
	rv := uint32(0)
	siz := uint32(unsafe.Sizeof(keepAlive))
	if err := syscall.WSAIoctl(syscall.Handle(s), syscall.SIO_KEEPALIVE_VALS, (*byte)(unsafe.Pointer(&keepAlive)), siz, nil, 0, &rv, nil, 0); err != nil {
		keepAlive.Time = prev
		return os.NewSyscallError("WSAIoctl", err)
	}
	return nil
}

func setKeepAliveProbeInterval(s uintptr, d time.Duration) error {
	keepAlive.Lock()
	defer keepAlive.Unlock()
	d += (time.Millisecond - time.Nanosecond)
	msecs := uint32(d / time.Millisecond)
	prev := keepAlive.Interval
	keepAlive.Interval = msecs
	rv := uint32(0)
	siz := uint32(unsafe.Sizeof(keepAlive))
	if err := syscall.WSAIoctl(syscall.Handle(s), syscall.SIO_KEEPALIVE_VALS, (*byte)(unsafe.Pointer(&keepAlive)), siz, nil, 0, &rv, nil, 0); err != nil {
		keepAlive.Interval = prev
		return os.NewSyscallError("WSAIoctl", err)
	}
	return nil
}

// See http://msdn.microsoft.com/en-us/library/windows/desktop/dd877220(v=vs.85).aspx
func setKeepAliveProbeCount(s uintptr, n int) error { return errOpNoSupport }
