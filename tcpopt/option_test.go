// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpopt_test

import (
	"testing"
	"time"

	"github.com/mikioh/tcp/tcpopt"
)

func TestOption(t *testing.T) {
	for _, o := range []tcpopt.Option{
		tcpopt.NoDelay(true),
		tcpopt.SendBuffer(1<<16 - 1),
		tcpopt.ReceiveBuffer(1<<16 - 1),
		tcpopt.KeepAlive(true),
		tcpopt.KeepAliveIdleInterval(1 * time.Hour),
		tcpopt.KeepAliveProbeInterval(10 * time.Minute),
		tcpopt.KeepAliveProbeCount(3),
		tcpopt.BufferUnsentThreshold(1<<16 - 1),
	} {
		if o.Level() == 0 {
			t.Fatalf("got %#x; want non-zero", o.Level())
		}
		if o.Name() == 0 {
			t.Fatalf("got %#x; want non-zero", o.Name())
		}
		if _, err := o.Marshal(); err != nil {
			t.Fatal(err)
		}
	}
}

const (
	testOptLevel = 0
	testOptName  = 1
)

type testOption struct{}

func (*testOption) Level() int                        { return testOptLevel }
func (*testOption) Name() int                         { return testOptName }
func (*testOption) Marshal() ([]byte, error)          { return make([]byte, 16), nil }
func parseTestOption(_ []byte) (tcpopt.Option, error) { return &testOption{}, nil }

func TestParse(t *testing.T) {
	var b [16]byte
	tcpopt.Register(testOptLevel, testOptName, parseTestOption)
	o, err := tcpopt.Parse(testOptLevel, testOptName, b[:])
	if err != nil {
		t.Fatal(err)
	}
	tcpopt.Unregister(testOptLevel, testOptName)
	o, err = tcpopt.Parse(testOptLevel, testOptName, b[:])
	if err == nil || o != nil {
		t.Fatalf("got %v, %v; want nil, error", o, err)
	}
}
