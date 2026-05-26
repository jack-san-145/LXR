package helper

import (
	"io"
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strconv"
)

func (h *Helper) SetupContainerNetworking(cb *models.ContainerBuilder) error {

	//use lxr defaulat bridge
	cb.Container.Bridge = h.NetworkConfig.UseBridge()

	//allocate ip address for container
	cb.Container.IpAddress = h.NetworkConfig.AllocateIp()

	log.Println("container ip: ", cb.Container.IpAddress)
	//create veth pairs
	cb.Container.ConVeth = cb.Container.ContainerName + "_Cveth" //one end that connected to container
	cb.Container.BrVeth = cb.Container.ContainerName + "_Bveth"  //other end that connected to bridge

	//create env variables
	container_pid_env := "CONTAINER_PID=" + strconv.Itoa(cb.Container.PID)
	container_ip_env := "CONTAINER_IP=" + cb.Container.IpAddress
	container_veth_env := "CONTAINER_VETH=" + cb.Container.ConVeth
	bridge_ip_env := "BRIDGE_IP=" + h.NetworkConfig.GetBrigeIp()
	bridge_veth_env := "BRIDGE_VETH=" + cb.Container.BrVeth

	cmd := exec.Command("bash", "../../script/ip-setup.sh")

	//inject env to the script
	cmd.Env = append(os.Environ(),
		container_pid_env,
		container_ip_env,
		container_veth_env,
		bridge_ip_env,
		bridge_veth_env,
	)

	// Get a pipe to read the command's standard output and error (stdout and stderr)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Stream stdoud(normal output) from the command to the client
	go func() {
		io.Copy(cb.Conn, stdout)
	}()

	// Stream stderr(error/warning output) from the command to the client
	go func() {
		io.Copy(cb.Conn, stderr)
	}()

	//start to run script at background
	err := cmd.Start()
	if err != nil {
		log.Println("container networking setup Error: ", err)
		return err
	}

	//wait until script to complete
	err = cmd.Wait()
	if err != nil {
		log.Println("container networking setup Error: ", err)
		return err

	}

	//find largest ip by comparing container ip and network's LastUsedIP
	largestIP := h.NetworkConfig.FindLargestIP(cb.Container.IpAddress, h.NetworkConfig.LastUsedIP)

	//update largest ip address to network's LastUsedIP
	h.NetworkConfig.SetLastUsedIp(largestIP)
	return nil

}
