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

	"github.com/mikioh/tcpopt"
)

// A Conn represents an end point that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	net.Conn
	s uintptr // socket descriptor for configuring options
}

func (c *Conn) setOption(level, name int, b []byte) error {
	s, err := socketOf(c.Conn)
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", setsockopt(s, level, name, b))
}

func (c *Conn) option(level, name int, b []byte) (int, error) {
	s, err := socketOf(c.Conn)
	if err != nil {
		return 0, err
	}
	n, err := getsockopt(s, level, name, b)
	return n, os.NewSyscallError("getsockopt", err)
}

func (c *Conn) buffered() int {
	s, err := socketOf(c.Conn)
	if err != nil {
		return -1
	}
	var b [4]byte
	if err := ioctl(s, options[soBuffered].name, b[:]); err != nil {
		return -1
	}
	return int(nativeEndian.Uint32(b[:]))
}

func (c *Conn) available() int {
	s, err := socketOf(c.Conn)
	if err != nil {
		return -1
	}
	var b [4]byte
	if runtime.GOOS == "darwin" {
		_, err = getsockopt(s, options[soAvailable].level, options[soAvailable].name, b[:])
	} else {
		err = ioctl(s, options[soAvailable].name, b[:])
	}
	if err != nil {
		return -1
	}
	n := int(nativeEndian.Uint32(b[:]))
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		var o tcpopt.SendBuffer
		_, err = getsockopt(s, o.Level(), o.Name(), b[:])
		if err != nil {
			return -1
		}
		n = int(nativeEndian.Uint32(b[:])) - n
	}
	return n
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

// NewConn returns a new end point.
func NewConn(c net.Conn) (*Conn, error) {
	s, err := socketOf(c)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: c, s: s}, nil
}
