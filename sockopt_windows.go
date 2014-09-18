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

var tcpKeepAlive = struct {
	sync.RWMutex
	syscall.TCPKeepalive
}{
	TCPKeepalive: syscall.TCPKeepalive{
		OnOff:    1,
		Time:     uint32(2 * time.Hour / time.Millisecond),
		Interval: uint32(time.Second / time.Millisecond),
	},
}

func (c *Conn) setKeepAliveIdleInterval(d time.Duration) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	tcpKeepAlive.Lock()
	defer tcpKeepAlive.Unlock()
	d += (time.Millisecond - time.Nanosecond)
	msecs := uint32(d / time.Millisecond)
	prev := tcpKeepAlive.Time
	tcpKeepAlive.Time = msecs
	rv := uint32(0)
	siz := uint32(unsafe.Sizeof(tcpKeepAlive))
	if err := syscall.WSAIoctl(fd, syscall.SIO_KEEPALIVE_VALS, (*byte)(unsafe.Pointer(&tcpKeepAlive)), siz, nil, 0, &rv, nil, 0); err != nil {
		tcpKeepAlive.Time = prev
		return os.NewSyscallError("WSAIoctl", err)
	}
	return nil
}

func (c *Conn) setKeepAliveProbeInterval(d time.Duration) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	tcpKeepAlive.Lock()
	defer tcpKeepAlive.Unlock()
	d += (time.Millisecond - time.Nanosecond)
	msecs := uint32(d / time.Millisecond)
	prev := tcpKeepAlive.Interval
	tcpKeepAlive.Interval = msecs
	rv := uint32(0)
	siz := uint32(unsafe.Sizeof(tcpKeepAlive))
	if err := syscall.WSAIoctl(fd, syscall.SIO_KEEPALIVE_VALS, (*byte)(unsafe.Pointer(&tcpKeepAlive)), siz, nil, 0, &rv, nil, 0); err != nil {
		tcpKeepAlive.Interval = prev
		return os.NewSyscallError("WSAIoctl", err)
	}
	return nil
}

func (c *Conn) setKeepAliveProbes(max int) error {
	// See http://msdn.microsoft.com/en-us/library/windows/desktop/dd877220(v=vs.85).aspx
	return errOpNoSupport
}

func (c *Conn) setCork(on bool) error {
	return errOpNoSupport
}

func (c *Conn) info() (*Info, error) {
	return nil, errOpNoSupport
}
