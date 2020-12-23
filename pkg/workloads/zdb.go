package workloads

type ZDB struct {
	Size            int64             `bson:"size" json:"size"`
	Mode            ZDBModeEnum       `bson:"mode" json:"mode"`
	Password        string            `bson:"password" json:"password"`
	DiskType        DiskTypeEnum      `bson:"disk_type" json:"disk_type"`
	Public          bool              `bson:"public" json:"public"`
	StatsAggregator []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
}

// WorkloadType implements Workload
func (z ZDB) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeZDB
}

type DiskTypeEnum uint8

const (
	DiskTypeHDD DiskTypeEnum = iota
	DiskTypeSSD
)

type ZDBModeEnum uint8

const (
	ZDBModeSeq ZDBModeEnum = iota
	ZDBModeUser
)
