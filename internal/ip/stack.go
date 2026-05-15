package ip

import ()

// stack to hold the ip addresses
type Stack struct {
	Stack []string
}

// ipstack to hold the ips and details of the network
type IpStack struct {
	Stack           Stack
	Network         string
	Cidr            int
	BridgeName      string
	BridgeIp        string
	IpStartRange    string
	IpEndRange      string
	NetworkAddr     string
	BroadCastAddr   string
	TotalUsableHost int
	HostUsed        int
	LastUsedIp      string
}

// function to create a new IpStack
func NewIpStack(ipstack IpStack) *IpStack {
	return &IpStack{
		Stack:           ipstack.Stack,
		Network:         ipstack.Network,
		Cidr:            ipstack.Cidr,
		BridgeName:      ipstack.BridgeName,
		BridgeIp:        ipstack.BridgeIp,
		IpStartRange:    ipstack.IpStartRange,
		IpEndRange:      ipstack.IpEndRange,
		NetworkAddr:     ipstack.NetworkAddr,
		BroadCastAddr:   ipstack.BroadCastAddr,
		TotalUsableHost: ipstack.TotalUsableHost,
		HostUsed:        ipstack.HostUsed,
		LastUsedIp:      ipstack.BridgeIp, //assign bridge ip to last used ip in this network
	}

}
