// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.9

package tcp

import (
	"errors"
	"net"
	"os"
	"reflect"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/mikioh/tcpopt"
)

// A Conn represents an end point that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	net.Conn
	s uintptr // socket descriptor for configuring options
}

func (c *Conn) ok() bool { return c != nil && c.Conn != nil }

func (c *Conn) setOption(level, name int, b []byte) error {
	return os.NewSyscallError("setsockopt", setsockopt(c.s, level, name, b))
}

func (c *Conn) option(level, name int, b []byte) (int, error) {
	n, err := getsockopt(c.s, level, name, b)
	return n, os.NewSyscallError("getsockopt", err)
}

func (c *Conn) buffered() int {
	var b [4]byte
	if err := ioctl(c.s, options[soBuffered].name, b[:]); err != nil {
		return -1
	}
	return int(nativeEndian.Uint32(b[:]))
}

func (c *Conn) available() int {
	var err error
	var b [4]byte
	if runtime.GOOS == "darwin" {
		_, err = getsockopt(c.s, options[soAvailable].level, options[soAvailable].name, b[:])
	} else {
		err = ioctl(c.s, options[soAvailable].name, b[:])
	}
	if err != nil {
		return -1
	}
	n := int(nativeEndian.Uint32(b[:]))
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		var o tcpopt.SendBuffer
		_, err = getsockopt(c.s, o.Level(), o.Name(), b[:])
		if err != nil {
			return -1
		}
		n = int(nativeEndian.Uint32(b[:])) - n
	}
	return n
}

// Close closes the connection.
func (c *Conn) Close() error {
	if !c.ok() {
		return syscall.EINVAL
	}
	atomic.StoreUintptr(&c.s, ^uintptr(0))
	return c.Conn.Close()
}

// NewConn returns a new end point.
func NewConn(c net.Conn) (*Conn, error) {
	s, err := socketOf(c)
	if err != nil {
		return nil, err
	}
	cc := &Conn{Conn: c}
	atomic.StoreUintptr(&cc.s, s)
	return cc, nil
}

func socketOf(c net.Conn) (uintptr, error) {
	switch c.(type) {
	case *net.TCPConn, *net.UDPConn, *net.IPConn:
		v := reflect.ValueOf(c)
		switch e := v.Elem(); e.Kind() {
		case reflect.Struct:
			fd := e.FieldByName("conn").FieldByName("fd")
			switch e := fd.Elem(); e.Kind() {
			case reflect.Struct:
				sysfd := e.FieldByName("sysfd")
				if runtime.GOOS == "windows" {
					return uintptr(sysfd.Uint()), nil
				}
				return uintptr(sysfd.Int()), nil
			}
		}
	}
	return 0, errors.New("invalid type")
}
