package ip

import (
	"fmt"
	"strconv"
	"strings"
)

// to allocate ip address and put next 100 ip addresses to stack
func (ips *IpStack) AllocateIp() {

	//ip addr to hold 100 ips temporarily
	var Ip_Arr []string

	//split all the four octet seperately
	net_arr := strings.Split(ips.LastUsedIp, ".")
	first, _ := strconv.Atoi(net_arr[0])
	second, _ := strconv.Atoi(net_arr[1])
	third, _ := strconv.Atoi(net_arr[2])
	fourth, _ := strconv.Atoi(net_arr[3])

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

// to create the ip with all 4 octets
func (ips *IpStack) MakeIp(first, second, third, fourth int) string {
	return fmt.Sprintf("%v.%v.%v.%v", first, second, third, fourth)
}

func (ips *IpStack) PushToStack(Ip_Arr []string) {

	//iterate ip array reversely to store the ips in asending order
	for i := len(Ip_Arr) - 1; i <= 0; i-- {

		ips.Stack.Push(Ip_Arr[i])
	}

}
