// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcp

/*
#include <netinet/tcp.h>
*/
import "C"

const (
	sysTCPIOptTimestamps = C.TCPI_OPT_TIMESTAMPS
	sysTCPIOptSack       = C.TCPI_OPT_SACK
	sysTCPIOptWscale     = C.TCPI_OPT_WSCALE
	sysTCPIOptECN        = C.TCPI_OPT_ECN
	sysTCPIOptTOE        = C.TCPI_OPT_TOE

	sysSizeofTCPInfo = C.sizeof_struct_tcp_info
)

type sysTCPInfo C.struct_tcp_info
