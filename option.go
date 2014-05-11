// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

// An Option represents a TCP option.
type Option interface {
	Kind() int
}

// A MaxSegSize represents a TCP maxiumum sengment size option.
type MaxSegSize int

// Kind returns a TCP option kind field.
func (mss MaxSegSize) Kind() int {
	return 2
}

// A WindowScale represents a TCP windows scale option.
type WindowScale int

// Kind returns a TCP option kind field.
func (ws WindowScale) Kind() int {
	return 3
}

// A SackPermitted represents a TCP selective acknowledgment permitted
// option.
type SackPermitted bool

// Kind returns a TCP option kind field.
func (sp SackPermitted) Kind() int {
	return 4
}
