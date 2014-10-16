// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import "time"

type sysSockoptLen int32

const (
	sysTCP_KEEPIDLE  = 0x4
	sysTCP_KEEPINTVL = 0x5
	sysTCP_KEEPCNT   = 0x6
	sysTCP_CORK      = 0x3
	sysTCP_INFO      = 0xb
)

var sockOpts = [ssoMax]sockOpt{
	ssoBuffered:               {sysSIOCINQ, ssoTypeInt, 0},
	ssoAvailable:              {sysSIOCOUTQ, ssoTypeInt, 0},
	ssoKeepAliveIdleInterval:  {sysTCP_KEEPIDLE, ssoTypeInt, time.Second},
	ssoKeepAliveProbeInterval: {sysTCP_KEEPINTVL, ssoTypeInt, time.Second},
	ssoKeepAliveProbes:        {sysTCP_KEEPCNT, ssoTypeInt, 0},
	ssoCork:                   {sysTCP_CORK, ssoTypeInt, 0},
	ssoInfo:                   {sysTCP_INFO, ssoTypeInfo, 0},
}
