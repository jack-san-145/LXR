package helper

import (
	"fmt"
	"log"
	"lxr-d/internal/models"
	"os"

	"strconv"
)

// to find equivalent bytes for given mb
func (h *Helper) MB(m int) string {
	return strconv.Itoa(m * 1024 * 1024) //mb -> kb -> byte
}

// to find percent for given integer
func (h *Helper) CPU(percent int) string {
	quota := percent * 1000
	return fmt.Sprintf("%v 100000", quota)
}

// to create new cgroup for container with limitations
func (h *Helper) CreateCgroup(con *models.Container) error {

	//get the container's init PID
	initPID, err := h.GetChildPID(con.PID)
	if err != nil {
		return err
	}

	base := "/sys/fs/cgroup/lxr/" + con.ContainerName

	//creates an base directory
	os.MkdirAll(base, 0755)

	//limits cpu,pid,ram
	err = os.WriteFile(base+"/cgroup.procs", []byte(initPID), 0644) //attach container init process with cgroup
	err = os.WriteFile(base+"/memory.max", []byte(h.MB(256)), 0644) //use 256mb of ram
	err = os.WriteFile(base+"/cpu.max", []byte(h.CPU(50)), 0644)    //use 50% of cpu
	err = os.WriteFile(base+"/pids.max", []byte("100"), 0644)       //can spawn upto 100 processes

	if err != nil {
		log.Println("cgroup err: ", err)
	}
	return err

}

// add nsenter's child to container's cgroup for resource control
func (h *Helper) AddNsenterChildToCgroup(nsenterPID int, containerName string) error {

	base := "/sys/fs/cgroup/lxr/" + containerName

	childPID, err := h.GetChildPID(nsenterPID)
	if err != nil {
		log.Println("child PID Error: ", err)
		return err
	}
	err = os.WriteFile(base+"/cgroup.procs", []byte(childPID), 0644) //attach nsenter child process with cgroup
	if err != nil {
		log.Println("cgroup attach Error: ", err)
		return err
	}
	return nil

}

// Freeze container by put 1 to cgroup.freeze
func (h *Helper) FreezeContainer(containerName string) error {

	base := "/sys/fs/cgroup/lxr/" + containerName

	// 1 = freeze
	err := os.WriteFile(base+"/cgroup.freeze", []byte("1"), 0644)

	if err != nil {
		log.Println("Freeze Error:", err)
		return err
	}

	return nil
}

// Unfreeze container
func (h *Helper) UnfreezeContainer(containerName string) error {

	base := "/sys/fs/cgroup/lxr/" + containerName

	// 0 = unfreeze
	err := os.WriteFile(base+"/cgroup.freeze", []byte("0"), 0644)

	if err != nil {
		log.Println("Unfreeze Error:", err)
		return err
	}

	return nil
}

// Kill all processes inside container cgroup
func (h *Helper) KillContainerCgroup(containerName string) error {

	base := "/sys/fs/cgroup/lxr/" + containerName

	// 1 = kill all processes in cgroup
	err := os.WriteFile(base+"/cgroup.kill", []byte("1"), 0644)

	if err != nil {
		log.Println("Kill Container cgroup Error:", err)
		return err
	}

	return nil
}
