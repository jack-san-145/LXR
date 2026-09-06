package helper

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

func (h *Helper) InitDaemon() (net.Listener, error) {

	lxr_sock := "/var/run/lxr.sock"
	err := os.Remove(lxr_sock)
	if err != nil {
		log.Println("Old Sock remove Error: ", err)
	}

	listener, err := net.Listen("unix", lxr_sock)
	if err != nil {
		log.Fatal("Listener Failed: ", err)
	}

	//changing sock path permissions
	group, err := user.LookupGroup("lxr")
	if err != nil {
		log.Println("Group not found: ", err)
	}
	group_id, _ := strconv.Atoi(group.Gid)

	err = os.Chown(lxr_sock, 0, group_id)
	if err != nil {
		log.Println("chown Error: ", err)
	}
	err = os.Chmod(lxr_sock, 0660)
	if err != nil {
		log.Println("chmod Error: ", err)
	}

	//run the lxr-init script
	cmd := exec.Command("../../script/lxr-init.sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println("LXR Initialization Failed..", err)
	}

	return listener, err
}
