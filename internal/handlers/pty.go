package handlers

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
)

func (h *LXRHandler) ExecHandler(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack not supported", http.StatusInternalServerError)
		return
	}
	conn, _, err := hijacker.Hijack()

	cmd := exec.Command(
		"nsenter",
		"--target", pid,
		"--pid", "--mount",
		"bash",
	)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Pty Error: ", err)
		return
	}
	defer ptmx.Close()

	go io.Copy(ptmx, conn)
	go io.Copy(conn, ptmx)

	cmd.Wait()
}
