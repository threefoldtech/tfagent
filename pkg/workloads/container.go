package workloads

import (
	"encoding/json"
	"net"
)

type Container struct {
	Flist             string              `bson:"flist" json:"flist"`
	HubUrl            string              `bson:"hub_url" json:"hub_url"`
	Environment       map[string]string   `bson:"environment" json:"environment"`
	SecretEnvironment map[string]string   `bson:"secret_environment" json:"secret_environment"`
	Entrypoint        string              `bson:"entrypoint" json:"entrypoint"`
	Interactive       bool                `bson:"interactive" json:"interactive"`
	Volumes           []ContainerMount    `bson:"volumes" json:"volumes"`
	NetworkConnection []NetworkConnection `bson:"network_connection" json:"network_connection"`
	Stats             []Stats             `bson:"stats" json:"stats"`
	Logs              []Logs              `bson:"logs" json:"logs"`
	Capacity          ContainerCapacity   `bson:"capcity" json:"capacity"`
}

type ContainerCapacity struct {
	Cpu      int64        `bson:"cpu" json:"cpu"`
	Memory   int64        `bson:"memory" json:"memory"`
	DiskSize uint64       `bson:"disk_size" json:"disk_size"`
	DiskType DiskTypeEnum `bson:"disk_type" json:"disk_type"`
}

type Logs struct {
	Type string    `bson:"type" json:"type"`
	Data LogsRedis `bson:"data" json:"data"`
}

type Stats struct {
	Type string          `bson:"type" json:"type"`
	Data json.RawMessage `bson:"data" json:"data"`
}

type LogsRedis struct {
	Stdout string `bson:"stdout" json:"stdout"`
	Stderr string `bson:"stderr" json:"stderr"`

	// Same as stdout, stderr urls but encrypted
	// with the node public key.
	SecretStdout string `bson:"secret_stdout" json:"secret_stdout"`
	SecretStderr string `bson:"secret_stderr" json:"secret_stderr"`
}

type StatsRedis struct {
	Endpoint string `bson:"endpoint" json:"endpoint"`
}

type ContainerMount struct {
	VolumeId   string `bson:"volume_id" json:"volume_id"`
	Mountpoint string `bson:"mountpoint" json:"mountpoint"`
}

type NetworkConnection struct {
	NetworkId   string `bson:"network_id" json:"network_id"`
	Ipaddress   net.IP `bson:"ipaddress" json:"ipaddress"`
	PublicIp6   bool   `bson:"public_ip6" json:"public_ip6"`
	YggdrasilIP bool   `bson:"yggdrasil_ip" json:"yggdrasil_ip"`
}
