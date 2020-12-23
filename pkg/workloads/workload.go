package workloads

import "encoding/json"

type (
	// WorkloadTypeEnum type of a workload
	WorkloadTypeEnum uint8

	// Workload definition
	Workload interface {
		WorkloadType() WorkloadTypeEnum
	}

	// WorkloadInfo generic info for a workload
	WorkloadInfo struct {
		Name     string
		Peer     string //TODO
		Workload Workload
	}

	// WorkloadEnveloppe info about a workload and its content
	WorkloadEnveloppe struct {
		Name         string
		Type         WorkloadTypeEnum
		WorkloadData json.RawMessage
	}
)

// available workload types
const (
	WorkloadTypeZDB WorkloadTypeEnum = iota
	WorkloadTypeContainer
	WorkloadTypeVolume
	WorkloadTypeNetwork
	WorkloadTypeKubernetes
	WorkloadTypeProxy
	WorkloadTypeReverseProxy
	WorkloadTypeSubDomain
	WorkloadTypeDomainDelegate
	WorkloadTypeGateway4To6
	WorkloadTypeNetworkResource
	WorkloadTypePublicIP
)
