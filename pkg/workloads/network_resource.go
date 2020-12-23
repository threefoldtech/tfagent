package workloads

import "net"

type NetworkResource struct {
	Name                         string            `bson:"name" json:"name"`
	NetworkIprange               net.IPNet         `bson:"network_iprange" json:"network_iprange"`
	WireguardPrivateKeyEncrypted string            `bson:"wireguard_private_key_encrypted" json:"wireguard_private_key_encrypted"`
	WireguardPublicKey           string            `bson:"wireguard_public_key" json:"wireguard_public_key"`
	WireguardListenPort          int64             `bson:"wireguard_listen_port" json:"wireguard_listen_port"`
	Iprange                      net.IPNet         `bson:"iprange" json:"iprange"`
	Peers                        []WireguardPeer   `bson:"peers" json:"peers"`
	StatsAggregator              []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
}

// WorkloadType implements Workload
func (z NetworkResource) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeNetworkResource
}

type WireguardPeer struct {
	PublicKey      string      `bson:"public_key" json:"public_key"`
	Endpoint       string      `bson:"endpoint" json:"endpoint"`
	Iprange        net.IPNet   `bson:"iprange" json:"iprange"`
	AllowedIprange []net.IPNet `bson:"allowed_iprange" json:"allowed_iprange"`
}
