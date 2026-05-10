package ip

import ()

type IpStack struct {
	Stack           *Stack
	Network         string
	Cidr            int
	IpStartRange    string
	IpEndRange      string
	NetworkAddr     string
	BroadCastAddr   string
	TotalUsableHost int
	HostUsed        int
	LastUsedIp      string
}

type Stack struct {
	Stack []string
}

func NewIpStack(ipstack IpStack) *IpStack {
	return &IpStack{
		Stack:           ipstack.Stack,
		Network:         ipstack.Network,
		Cidr:            ipstack.Cidr,
		IpStartRange:    ipstack.IpStartRange,
		IpEndRange:      ipstack.IpEndRange,
		NetworkAddr:     ipstack.NetworkAddr,
		BroadCastAddr:   ipstack.BroadCastAddr,
		TotalUsableHost: ipstack.TotalUsableHost,
		HostUsed:        ipstack.HostUsed,
		LastUsedIp:      ipstack.NetworkAddr,
	}

}
