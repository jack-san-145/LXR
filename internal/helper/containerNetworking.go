package helper

import (
	"fmt"
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strconv"
)

func (h *Helper) SetupContainerNetworking(con *models.Container) error {

	//use lxr defaulat bridge
	con.Bridge = h.NetworkConfig.UseBridge()

	//allocate ip address for container
	con.IpAddress = h.NetworkConfig.AllocateIp()

	log.Println("container ip: ", con.IpAddress)
	//create veth pairs
	con.ConVeth = con.ContainerName + "_Cveth" //one end that connected to container
	con.BrVeth = con.ContainerName + "_Bveth"  //other end that connected to bridge

	//create env variables
	container_pid_env := "CONTAINER_PID=" + strconv.Itoa(con.PID)
	container_ip_env := "CONTAINER_IP=" + con.IpAddress
	container_veth_env := "CONTAINER_VETH=" + con.ConVeth
	bridge_ip_env := "BRIDGE_IP=" + h.NetworkConfig.GetBrigeIp()
	bridge_veth_env := "BRIDGE_VETH=" + con.BrVeth

	cmd := exec.Command("bash", "../../script/ip-setup.sh")

	//inject env to the script
	cmd.Env = append(os.Environ(),
		container_pid_env,
		container_ip_env,
		container_veth_env,
		bridge_ip_env,
		bridge_veth_env,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//run script in foreground
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error container networking setup : ", err)
		return err
	}

	//update last used ip address with container ip
	h.NetworkConfig.SetLastUsedIp(con.IpAddress)
	return nil
}
