package helper

import (
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strconv"
)

func (h *Helper) StopContainer(name string) (bool, error) {

	pid, exists := h.GetContainerPid(name)
	if !exists {
		return exists, nil
	}
	pid_int, _ := strconv.Atoi(pid)
	ps, err := os.FindProcess(pid_int)

	if err != nil {
		return exists, err
	}

	err = ps.Kill()
	if err == nil {
		delete(h.ContainerManager.ActiveContainers, name)
	}
	return exists, err
}

func (h *Helper) KillContainer(con *models.Container) error {

	container_name_env := "CONTAINER_NAME=" + con.ContainerName
	container_id_env := "CONTAINER_ID=" + con.ContainerId
	container_pid_env := "CONTAINER_PID=" + strconv.Itoa(con.PID)
	bridge_veth_env := "BRIDGE_VETH=" + con.BrVeth

	cmd := exec.Command("bash", "../../script/container-kill.sh")

	//inject env to the script
	cmd.Env = append(os.Environ(),
		container_name_env,
		container_id_env,
		bridge_veth_env,
		container_pid_env,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//run script in foreground
	err := cmd.Run()
	if err != nil {
		log.Println("Error container kill: ", err)
		return err
	}

	//remove killed container from containerManager
	delete(h.ContainerManager.ActiveContainers, con.ContainerName)
	delete(h.ContainerManager.AllContainers, con.ContainerName)

	//put back the container's ip address to IPPool
	h.NetworkConfig.ReturnIPToPool(con.IpAddress)

	//kill container cgroup if exists
	exists := h.CgroupExists(con.ContainerName)
	if exists {
		err := h.KillContainerCgroup(con.ContainerName)
		return err

	}
	return nil
}
