package workloads

type Volume struct {
	Size int64          `bson:"size" json:"size"`
	Type VolumeTypeEnum `bson:"type" json:"type"`
}

// WorkloadType implements Workload
func (z Volume) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeVolume
}

type VolumeTypeEnum uint8

const (
	VolumeTypeHDD VolumeTypeEnum = iota
	VolumeTypeSSD
)
