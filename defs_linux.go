// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>

#include <linux/sockios.h>

#include <netinet/tcp.h>
*/
import "C"

const (
	sysSIOCINQ  = C.SIOCINQ
	sysSIOCOUTQ = C.SIOCOUTQ

	sysTCPI_OPT_TIMESTAMPS = C.TCPI_OPT_TIMESTAMPS
	sysTCPI_OPT_SACK       = C.TCPI_OPT_SACK
	sysTCPI_OPT_WSCALE     = C.TCPI_OPT_WSCALE
	sysTCPI_OPT_ECN        = C.TCPI_OPT_ECN
	sysTCPI_OPT_ECN_SEEN   = C.TCPI_OPT_ECN_SEEN
	sysTCPI_OPT_SYN_DATA   = C.TCPI_OPT_SYN_DATA

	CAOpen     CAState = C.TCP_CA_Open
	CADisorder CAState = C.TCP_CA_Disorder
	CACWR      CAState = C.TCP_CA_CWR
	CARecovery CAState = C.TCP_CA_Recovery
	CALoss     CAState = C.TCP_CA_Loss

	sysSizeofTCPInfo = C.sizeof_struct_tcp_info
)

type sysTCPInfo C.struct_tcp_info
