package helpers

import (
	"log"
	"os/exec"
	"syscall"
	"time"
)

// Run bash command in subprocess with tiemout contorl.
func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var err error

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		log.Println("timeout and killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

		//IMPORTANT: cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is necessary before cmd.Start()
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			log.Println("kill failed", err)
		}

		return err, true
	case err = <-done:
		return err, false
	}
}
