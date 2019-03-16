// Collect Intel GPU render busy info via package intel-gpu-tools.
//     apt-get install intel-gpu-tools
package GPUIntel

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/OpenCCTV/gpu/gpu/helpers"
)

const (
	pathBin = "/usr/bin/intel_gpu_top"
)

func Gets(debug bool) (m *[]map[string]interface{}, err error) {
	var errmsg string
	result := []map[string]interface{}{}
	timeoutInSeconds := 2

	if _, err = os.Stat(pathBin); os.IsNotExist(err) {
		errmsg = "stat tool not found"
		if debug {
			log.Println(errmsg, pathBin)
		}
		err = errors.New(errmsg)
		return
	}

	tmpfile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
	defer os.Remove(tmpfile.Name())

	arg := fmt.Sprintf("sudo timeout %d %s -o %s", timeoutInSeconds, pathBin, tmpfile.Name())
	cmd := exec.Command("bash", "-c", arg)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	err, isTimeout := helpers.CmdRunWithTimeout(cmd, time.Duration(timeoutInSeconds+2)*time.Second)
	// see also https://stackoverflow.com/questions/10385551/get-exit-code-go
	TIMEOUT_EXIT := 124
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitStatus := status.ExitStatus()
				if exitStatus != TIMEOUT_EXIT {
					if debug {
						log.Println(err)
					}
					return
				}
			}
		} else {
			if debug {
				log.Println(err)
			}
			return
		}
	}

	if isTimeout {
		errmsg = "run cmd timeout"
		if debug {
			log.Println(errmsg)
		}
		err = errors.New(errmsg)
		return
	}

	stderrStr := strings.TrimSpace(stderr.String())
	if stderrStr != "" {
		if debug {
			log.Println(stderrStr)
		}
		err = errors.New(stderrStr)
		return
	}

	resultStat, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}

	lines := strings.Split(string(resultStat), "\n")
	if len(lines) < 2 {
		errmsg = "got unexpected stat tool output"
		if debug {
			log.Println(errmsg)
		}
		err = errors.New(errmsg)
		return

	}
	columns := strings.Fields(lines[1])
	if len(columns) < 2 {
		errmsg = "parse stat tool output failed"
		if debug {
			log.Println(errmsg)
		}
		err = errors.New(errmsg)
		return
	}

	item := map[string]interface{}{}
	item["renderbusy"] = columns[1]
	// NOTICE: hard-coded its GPUID
	item["gpuid"] = "0"
	item["vendor"] = "intel"
	result = append(result, item)

	return &result, nil
}
