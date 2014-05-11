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
	sysTCPIOptECNSeen    = C.TCPI_OPT_ECN_SEEN
	sysTCPIOptSynData    = C.TCPI_OPT_SYN_DATA

	CAOpen     CAState = C.TCP_CA_Open
	CADisorder CAState = C.TCP_CA_Disorder
	CACWR      CAState = C.TCP_CA_CWR
	CARecovery CAState = C.TCP_CA_Recovery
	CALoss     CAState = C.TCP_CA_Loss
)
