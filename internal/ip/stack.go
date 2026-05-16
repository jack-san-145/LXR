package ip

import ()

// stack to hold the ip addresses
type Stack struct {
	IPPool []string
}

// NetworkConfig to hold the ipstack and details of the network
type NetworkConfig struct {
	IPStack         Stack
	Network         string
	CIDR            int
	BridgeName      string
	BridgeIP        string
	IPStartRange    string
	IPEndRange      string
	NetworkAddr     string
	BroadcastAddr   string
	TotalUsableHost int
	HostUsed        int
	LastUsedIP      string
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
