// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <sys/ioctl.h>

#include <linux/inet_diag.h>
#include <linux/sockios.h>
#include <linux/tcp.h>
*/
import "C"

const (
	sysSIOCINQ  = C.SIOCINQ
	sysSIOCOUTQ = C.SIOCOUTQ

	sysTCP_INFO       = C.TCP_INFO
	sysTCP_CONGESTION = C.TCP_CONGESTION
	sysTCP_CC_INFO    = C.TCP_CC_INFO

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

	sizeofTCPInfo      = C.sizeof_struct_tcp_info
	sizeofTCPCCInfo    = C.sizeof_union_tcp_cc_info
	sizeofTCPVegasInfo = C.sizeof_struct_tcpvegas_info
	sizeofTCPDCTCPInfo = C.sizeof_struct_tcp_dctcp_info
)

type sysTCPInfo C.struct_tcp_info

type sysTCPCCInfo C.union_tcp_cc_info

type sysTCPVegasInfo C.struct_tcpvegas_info

type sysTCPDCTCPInfo C.struct_tcp_dctcp_info
