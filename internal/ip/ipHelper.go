package ip

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

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

func (nc *NetworkConfig) FindLargestIP(IP1, IP2 string) string {

	//split both ip's as octet array
	IP1_arr := strings.Split(IP1, ".")
	IP2_arr := strings.Split(IP2, ".")

	//traverse ip array
	for i := 0; i < len(IP1_arr); i++ {

		//convert string octet to int
		ip1_octet, _ := strconv.Atoi(IP1_arr[i])
		ip2_octet, _ := strconv.Atoi(IP2_arr[i])

		//return largest ip as ip1 if ip1 octet is greater
		if ip1_octet > ip2_octet {
			return IP1
		} else if ip2_octet > ip1_octet { //return largest ip as ip2 if ip2 octet is greater
			return IP2
		}
		//if both actet are same do nothing, simply continue iteration
	}
	//if for IPs have same octets,then return any one ip from it(both are same)
	return IP1
}
