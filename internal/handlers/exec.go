package handlers

import (
	"io"
	"log"
	"lxr-d/internal/response"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
)

func (h *Handler) ExecHandler(w http.ResponseWriter, r *http.Request) {

	//check the container is running if true get its pid
	con_name := r.URL.Query().Get("container_name")
	pid, ok := h.Helper.GetContainerPid(con_name)
	if !ok || pid == "" {
		log.Println("container not running")
		response.WriteJson(w, "Container not Found\n")
	}

	//hijack the http connnection and stream i/o in real-time over uds
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack not supported", http.StatusInternalServerError)
		return
	}
	conn, _, err := hijacker.Hijack()

	//run cmd in pseudo terminal for given container
	cmd := exec.Command(
		"nsenter",
		"--target", pid,
		"--pid", "--mount", "--uts",
		"bash",
	)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Pty Error: ", err)
		return
	}

	defer func() {
		ptmx.Close()
		conn.Close()
	}()

	go io.Copy(ptmx, conn)
	go io.Copy(conn, ptmx)

	cmd.Wait()
}
