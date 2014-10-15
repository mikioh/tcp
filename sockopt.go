// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "time"

// Socket options
const (
	ssoKeepAliveIdleInterval = iota
	ssoKeepAliveProbeInterval
	ssoKeepAliveProbes
	ssoCork
	ssoInfo
	ssoMax
)

// Socket option value types
const (
	ssoTypeInt = iota + 1
	ssoTypeInfo
)

// A sockOpt represents a binding for socket option.
type sockOpt struct {
	name int           // option name, must be equal or greater than 1
	typ  int           // option value type, must be equal or greater than 1
	tmu  time.Duration // unit of time
}
