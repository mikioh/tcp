// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>

#include <linux/sockios.h>
*/
import "C"

const (
	sysSIOCINQ  = C.SIOCINQ
	sysSIOCOUTQ = C.SIOCOUTQ
)
