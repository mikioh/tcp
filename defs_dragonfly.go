// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>

#include <netinet/tcp.h>
*/
import "C"

const (
	sysFIONREAD = C.FIONREAD

	sysTCP_NOPUSH    = C.TCP_NOPUSH
	sysTCP_KEEPIDLE  = C.TCP_KEEPIDLE
	sysTCP_KEEPINTVL = C.TCP_KEEPINTVL
	sysTCP_KEEPCNT   = C.TCP_KEEPCNT
)
