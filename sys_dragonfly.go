// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly

package tcp

import "time"

type sysSockoptLen int32

var sockOpts = [ssoMax]sockOpt{
	ssoBuffered:               {sysFIONREAD, ssoTypeInt, 0},
	ssoCork:                   {sysTCP_NOPUSH, ssoTypeInt, 0},
	ssoKeepAliveIdleInterval:  {sysTCP_KEEPIDLE, ssoTypeInt, time.Millisecond},
	ssoKeepAliveProbeInterval: {sysTCP_KEEPINTVL, ssoTypeInt, time.Millisecond},
	ssoKeepAliveProbeCount:    {sysTCP_KEEPCNT, ssoTypeInt, 0},
}
