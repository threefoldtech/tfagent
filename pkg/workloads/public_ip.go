package workloads

import (
	"net"
)

// PublicIP is a struct that defines the workload to reserve a public ip on the grid
type PublicIP struct {
	IPaddress net.IPNet `bson:"ipaddress" json:"ipaddress"`
}

// WorkloadType implements Workload
func (z PublicIP) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypePublicIP
}
