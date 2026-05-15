package ip

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// return lxr's default bridge 'lxr0'
func (ips *IpStack) UseBridge() string {
	return ips.BridgeName
}

func (ips *IpStack) GetBrigeIp() string {
	log.Println("bridge ip from get : ", ips.BridgeIp)
	return ips.BridgeIp
}

func (ips *IpStack) SetLastUsedIp(ip string) {
	ips.LastUsedIp = ip
}

// to create the ip with all 4 octets
func (ips *IpStack) MakeIp(first, second, third, fourth int) string {
	return fmt.Sprintf("%v.%v.%v.%v/%v", first, second, third, fourth, ips.Cidr)
}

func (ips *IpStack) PushToStack(Ip_Arr []string) {

	var count int
	//iterate ip array reversely to store the ips in asending order
	for i := len(Ip_Arr) - 1; i >= 0; i-- {

		ips.Stack.Push(Ip_Arr[i])
		count++
	}

}

func (ips *IpStack) AllocateIp() string {
	if ips.Stack.IsEmpty() {
		ips.RefillIp()
	}
	ip, err := ips.Stack.Pop()
	if err != nil {
		log.Println("pop error:", err)
	}
	log.Println("allocated ip in stack: ", ip)
	return ip

}

// to allocate ip address and put next 100 ip addresses to stack
func (ips *IpStack) RefillIp() {

	//ip addr to hold 100 ips temporarily
	var Ip_Arr []string

	//split all the four octet seperately
	net_arr := strings.Split(ips.LastUsedIp, ".")
	first, _ := strconv.Atoi(net_arr[0])
	second, _ := strconv.Atoi(net_arr[1])
	third, _ := strconv.Atoi(net_arr[2])

	//extract 10.10.10.10 from 10.10.10.10/17
	fourth, _ := strconv.Atoi(strings.Split(net_arr[3], "/")[0])

	//find broadcast addr 3rd octet
	broadcast_addr, _ := strconv.Atoi(strings.Split(ips.BroadCastAddr, ".")[2])

	//loop to allocate next 100 ips
	for i := 1; i <= 100; i++ {

		//ensure octet not exceed limit '255'
		if fourth < 255 {
			fourth++

			//this avoid to take broadcast address as an host ip
			if third == broadcast_addr && fourth == 255 {
				break
			}

			//now form the ip with all 4 octets
			ip := ips.MakeIp(first, second, third, fourth)
			Ip_Arr = append(Ip_Arr, ip)

		} else {
			//to increment the 3rd octet if limit '255' reached
			if third < broadcast_addr {
				third++
				fourth = -1 //assign -1 to take next ip from 0
			}

		}
	}
	ips.PushToStack(Ip_Arr)

}
