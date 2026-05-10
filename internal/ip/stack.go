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
