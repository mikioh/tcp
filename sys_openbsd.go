// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

const sysTCP_NOPUSH = 0x10

var sockOpts = [ssoMax]sockOpt{
	ssoCork: {sysTCP_NOPUSH, ssoTypeInt, 0},
}
