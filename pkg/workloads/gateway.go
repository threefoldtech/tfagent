package workloads

type GatewayProxy struct {
	Domain  string `bson:"domain" json:"domain"`
	Addr    string `bson:"addr" json:"addr"`
	Port    uint32 `bson:"port" json:"port"`
	PortTLS uint32 `bson:"port_tls" json:"port_tls"`
}

// WorkloadType implements Workload
func (z GatewayProxy) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeProxy
}

type GatewayReverseProxy struct {
	Domain string `bson:"domain" json:"domain"`
	Secret string `bson:"secret" json:"secret"`
}

// WorkloadType implements Workload
func (z GatewayReverseProxy) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeReverseProxy
}

type GatewaySubdomain struct {
	Domain string   `bson:"domain" json:"domain"`
	IPs    []string `bson:"ips" json:"ips"`
}

// WorkloadType implements Workload
func (z GatewaySubdomain) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeSubDomain
}

type GatewayDelegate struct {
	Domain string `bson:"domain" json:"domain"`
}

// WorkloadType implements Workload
func (z GatewayDelegate) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeDomainDelegate
}

type Gateway4To6 struct {
	PublicKey string `bson:"public_key" json:"public_key"`
}

// WorkloadType implements Workload
func (z Gateway4To6) WorkloadType() WorkloadTypeEnum {
	return WorkloadTypeGateway4To6
}
