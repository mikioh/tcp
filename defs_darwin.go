// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>
#include <sys/socket.h>
*/
import "C"

const (
	sysSOL_SOCKET = C.SOL_SOCKET

	sysFIONREAD = C.FIONREAD

	sysSO_NREAD     = C.SO_NREAD
	sysSO_NWRITE    = C.SO_NWRITE
	sysSO_NUMRCVPKT = C.SO_NUMRCVPKT
)
