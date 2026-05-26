package helper

import (
	"fmt"
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

// to get container init PID
func (h *Helper) GetContainerInitPid(unsharePID int) (string, error) {

	cmd := exec.Command(
		"pgrep",
		"-P",
		strconv.Itoa(unsharePID),
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// to create new cgroup for container with limitations
func (h *Helper) CreateCgroup(con *models.Container) error {

	initPID, err := h.GetContainerInitPid(con.PID)
	if err != nil {
		return err
	}

	base := "/sys/fs/cgroup/lxr/" + con.ContainerName + con.ContainerId

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
