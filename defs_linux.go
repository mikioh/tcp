// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

// +godefs map struct_in_addr [4]byte /* in_addr */
// +godefs map struct_in6_addr [16]byte /* in6_addr */

/*
#include <sys/ioctl.h>
#include <sys/socket.h>

#include <linux/if.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/netfilter_ipv4.h>
#include <linux/netfilter_ipv6/ip6_tables.h>
#include <linux/sockios.h>
*/
import "C"

const (
	sysSO_ORIGINAL_DST      = C.SO_ORIGINAL_DST
	sysIP6T_SO_ORIGINAL_DST = C.IP6T_SO_ORIGINAL_DST
)

type sockaddrStorage C.struct_sockaddr_storage

type sockaddr C.struct_sockaddr

type sockaddrInet C.struct_sockaddr_in

type sockaddrInet6 C.struct_sockaddr_in6

const (
	sizeofSockaddrStorage = C.sizeof_struct_sockaddr_storage
	sizeofSockaddr        = C.sizeof_struct_sockaddr
	sizeofSockaddrInet    = C.sizeof_struct_sockaddr_in
	sizeofSockaddrInet6   = C.sizeof_struct_sockaddr_in6
)
