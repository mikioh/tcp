// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.3

package tcp

import "time"

type sysSockoptLen int32

const (
	sysTCP_KEEPIDLE  = 0x22
	sysTCP_KEEPINTVL = 0x24
	sysTCP_KEEPCNT   = 0x23
	sysTCP_CORK      = 0x18
)

var sockOpts = [ssoMax]sockOpt{
	ssoKeepAliveIdleInterval:  {sysTCP_KEEPIDLE, ssoTypeInt, time.Second},
	ssoKeepAliveProbeInterval: {sysTCP_KEEPINTVL, ssoTypeInt, time.Second},
	ssoKeepAliveProbes:        {sysTCP_KEEPCNT, ssoTypeInt, 0},
	ssoCork:                   {sysTCP_CORK, ssoTypeInt, 0},
}
