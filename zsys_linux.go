// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs sys_linux.go

package tcp

const (
	sysTCPIOptTimestamps = 0x1
	sysTCPIOptSack       = 0x2
	sysTCPIOptWscale     = 0x4
	sysTCPIOptECN        = 0x8
	sysTCPIOptECNSeen    = 0x10
	sysTCPIOptSynData    = 0x20

	CAOpen     CAState = 0x0
	CADisorder CAState = 0x1
	CACWR      CAState = 0x2
	CARecovery CAState = 0x3
	CALoss     CAState = 0x4
)
