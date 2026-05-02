package helper

import (
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
)

func (h *Helper) InitDaemon() net.Listener {

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

	return listener
}
