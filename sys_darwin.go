// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "time"

type sysSockoptLen int32

const (
	sysTCP_KEEPALIVE = 0x10
	sysTCP_KEEPINTVL = 0x101
	sysTCP_KEEPCNT   = 0x102
	sysTCP_NOPUSH    = 0x4
)

var sockOpts = [ssoMax]sockOpt{
	ssoReadBufferLen:          {sysFIONREAD, ssoTypeInt, 0},
	ssoKeepAliveIdleInterval:  {sysTCP_KEEPALIVE, ssoTypeInt, time.Second},
	ssoKeepAliveProbeInterval: {sysTCP_KEEPINTVL, ssoTypeInt, time.Second},
	ssoKeepAliveProbes:        {sysTCP_KEEPCNT, ssoTypeInt, 0},
	ssoCork:                   {sysTCP_NOPUSH, ssoTypeInt, 0},
}
