// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

// +godefs map struct_in_addr [4]byte /* in_addr */
// +godefs map struct_in6_addr [16]byte /* in6_addr */
// +godefs map struct_pf_addr [16]byte /* pf_addr */

/*
#include <sys/ioctl.h>
#include <sys/socket.h>

#include <netinet/in.h>

#include <net/if.h>
#include <net/pfvar.h>
*/
import "C"

const (
	sysFIONREAD = C.FIONREAD

	sysAF_INET  = C.AF_INET
	sysAF_INET6 = C.AF_INET6

	sysPF_INOUT = C.PF_INOUT
	sysPF_IN    = C.PF_IN
	sysPF_OUT   = C.PF_OUT
	sysPF_FWD   = C.PF_FWD

	sysDIOCNATLOOK = C.DIOCNATLOOK
)

type sockaddrStorage C.struct_sockaddr_storage

type sockaddr C.struct_sockaddr

type sockaddrInet C.struct_sockaddr_in

type sockaddrInet6 C.struct_sockaddr_in6

type pfiocNatlook C.struct_pfioc_natlook

const (
	sizeofSockaddrStorage = C.sizeof_struct_sockaddr_storage
	sizeofSockaddr        = C.sizeof_struct_sockaddr
	sizeofSockaddrInet    = C.sizeof_struct_sockaddr_in
	sizeofSockaddrInet6   = C.sizeof_struct_sockaddr_in6
	sizeofPfiocNatlook    = C.sizeof_struct_pfioc_natlook
)
