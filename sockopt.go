// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "time"

// Sticky socket options
const (
	ssoBuffered = iota
	ssoAvailable
	ssoCork
	ssoUnsentThreshold
	ssoKeepAliveIdleInterval
	ssoKeepAliveProbeInterval
	ssoKeepAliveProbeCount
	ssoMax
)

// Sticky socket option value types
const (
	ssoTypeInt = iota + 1
)

// A sockOpt represents a binding for sticky socket option.
type sockOpt struct {
	name int           // option name, must be equal or greater than 1
	typ  int           // option value type, must be equal or greater than 1
	tmu  time.Duration // unit of time
}
