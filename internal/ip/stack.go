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

func (s Stack) Size() int {
	return len(s.Stack)
}

func (s Stack) IsEmpty() bool {
	if s.Size() != 0 {
		return false
	}
	return true
}

func (s Stack) Push(ip string) {
	s.Stack = append(s.Stack, ip)
}
