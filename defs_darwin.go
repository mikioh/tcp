// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// +godefs map struct_in_addr [4]byte /* in_addr */
// +godefs map struct_in6_addr [16]byte /* in6_addr */
// +godefs map struct_pf_addr [16]byte /* pf_addr */

package tcp

/*
#include <sys/ioctl.h>
#include <sys/socket.h>

#include <netinet/in.h>

union pf_state_xport {
	u_int16_t port;
	u_int16_t call_id;
	u_int32_t spi;
};

struct pf_addr {
	uint8_t addr8[16];
};

struct pfioc_natlook {
	struct pf_addr saddr;
	struct pf_addr daddr;
	struct pf_addr rsaddr;
	struct pf_addr rdaddr;
	union pf_state_xport sxport;
	union pf_state_xport dxport;
	union pf_state_xport rsxport;
	union pf_state_xport rdxport;
	sa_family_t af;
	u_int8_t proto;
	u_int8_t proto_variant;
	u_int8_t direction;
};

#define DIOCNATLOOK _IOWR('D', 23, struct pfioc_natlook)
*/
import "C"

const (
	sysSOL_SOCKET = C.SOL_SOCKET

	sysFIONREAD = C.FIONREAD

	sysSO_NREAD     = C.SO_NREAD
	sysSO_NWRITE    = C.SO_NWRITE
	sysSO_NUMRCVPKT = C.SO_NUMRCVPKT

	sysAF_INET  = C.AF_INET
	sysAF_INET6 = C.AF_INET6

	sysPF_INOUT = 0
	sysPF_IN    = 1
	sysPF_OUT   = 2

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
