// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

// An OptionKind represents a TCP option kind.
type OptionKind int

const (
	KindMaxSegSize    OptionKind = 2
	KindWindowScale   OptionKind = 3
	KindSackPermitted OptionKind = 4
	KindTimestamps    OptionKind = 8
)

var optionKinds = map[OptionKind]string{
	KindMaxSegSize:    "mss",
	KindWindowScale:   "wscale",
	KindSackPermitted: "sack perm",
	KindTimestamps:    "tmstamps",
}

func (k OptionKind) String() string {
	s, ok := optionKinds[k]
	if !ok {
		return "<nil>"
	}
	return s
}

// An Option represents a TCP option.
type Option interface {
	Kind() OptionKind
}

// A MaxSegSize represents a TCP maxiumum sengment size option.
type MaxSegSize uint

// Kind returns a TCP option kind field.
func (mss MaxSegSize) Kind() OptionKind {
	return KindMaxSegSize
}

// A WindowScale represents a TCP windows scale option.
type WindowScale int

// Kind returns a TCP option kind field.
func (ws WindowScale) Kind() OptionKind {
	return KindWindowScale
}

// A SackPermitted reports whether a TCP selective acknowledgment
// permitted option is enabled.
type SackPermitted bool

// Kind returns a TCP option kind field.
func (sp SackPermitted) Kind() OptionKind {
	return KindSackPermitted
}

// A Timestamps reports whether a TCP timestamps option is enabled.
type Timestamps bool

// Kind returns a TCP option kind field.
func (ts Timestamps) Kind() OptionKind {
	return KindTimestamps
}
