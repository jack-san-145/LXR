package ip

import ()

// stack to hold the ip addresses
type Stack struct {
	IPPool []string `json:"ip_pool"`
}

// NetworkConfig to hold the ipstack and details of the network
type NetworkConfig struct {
	IPStack         Stack  `json:"ip_stack"`
	Network         string `json:"network"`
	CIDR            int    `json:"cidr"`
	BridgeName      string `json:"bridge_name"`
	BridgeIP        string `json:"bridge_ip"`
	IPStartRange    string `json:"ip_start_range"`
	IPEndRange      string `json:"ip_end_range"`
	NetworkAddr     string `json:"network_addr"`
	BroadcastAddr   string `json:"broadcast_addr"`
	TotalUsableHost int    `json:"total_usable_host"`
	HostUsed        int    `json:"host_used"`
	LastUsedIP      string `json:"last_used_ip"`
}

// function to create a new NetworkConfig
func NewNetwork(cfg NetworkConfig) *NetworkConfig {
	return &NetworkConfig{
		IPStack:         cfg.IPStack,
		Network:         cfg.Network,
		CIDR:            cfg.CIDR,
		BridgeName:      cfg.BridgeName,
		BridgeIP:        cfg.BridgeIP,
		IPStartRange:    cfg.IPStartRange,
		IPEndRange:      cfg.IPEndRange,
		NetworkAddr:     cfg.NetworkAddr,
		BroadcastAddr:   cfg.BroadcastAddr,
		TotalUsableHost: cfg.TotalUsableHost,
		HostUsed:        cfg.HostUsed,
		LastUsedIP:      cfg.BridgeIP, //assign bridge ip to last used ip in this network
	}

}
