package ip

import ("fmt"
"log")

// return lxr's default bridge 'lxr0'
func (nc *NetworkConfig) UseBridge() string {
	return nc.BridgeName
}

func (nc *NetworkConfig) GetBrigeIp() string {
	log.Println("bridge ip from get : ", nc.BridgeIP)
	return nc.BridgeIP
}

func (nc *NetworkConfig) SetLastUsedIp(ip string) {
	nc.LastUsedIP = ip
}

// to create the ip with all 4 octets
func (nc *NetworkConfig) MakeIp(first, second, third, fourth int) string {
	return fmt.Sprintf("%v.%v.%v.%v/%v", first, second, third, fourth, nc.CIDR)
}

