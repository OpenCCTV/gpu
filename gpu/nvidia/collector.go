// Collect NVIDIA GPU info via package nvidia-381.
//    add-apt-repository ppa:graphics-drivers/ppa
//    apt-get install nvidia-384
package GPUNvidia

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/OpenCCTV/gpu/gpu/helpers"
)

const (
	pathBin = "/usr/bin/nvidia-smi"
)

func Gets(debug bool) (m *[]map[string]interface{}, err error) {
	result := []map[string]interface{}{}

	total, err := getTotalGPU(debug)
	if err != nil {
		return
	}

	cmd := exec.Command("bash", "-c", "/usr/bin/nvidia-smi -pm 1 > /dev/null &&"+pathBin)
	out, err := cmd.Output()
	if err != nil {
		return
	}

	output := strings.TrimSpace(string(out))
	lines := strings.Split(output, "\n")

	FIRST_ID0_GPU_ROW := 8
	ROWS_OFFSET := 3

	if len(lines) >= (FIRST_ID0_GPU_ROW + ROWS_OFFSET*total) {

		i := FIRST_ID0_GPU_ROW
		for i < FIRST_ID0_GPU_ROW+ROWS_OFFSET*total {
			gpuID := (i - FIRST_ID0_GPU_ROW) / ROWS_OFFSET
			part, err := ParseGPURow(lines[i])
			if err != nil {
				if debug {
					log.Println("parse stat info failed, GPU ID", gpuID)
				}
			} else {
				(*part)["gpuid"] = gpuID
				(*part)["vendor"] = "nvidia"
				result = append(result, *part)
			}

			i += ROWS_OFFSET
		}
	}

	return &result, nil
}

func getTotalGPU(debug bool) (total int, err error) {
	var errmsg string

	timeoutInSeconds := 1

	if _, err = os.Stat(pathBin); os.IsNotExist(err) {
		errmsg = "stat tool not found"
		if debug {
			log.Println(errmsg, pathBin)
		}
		err = errors.New(errmsg)
		return
	}

	arg := fmt.Sprintf(`%s --list-gpus`, pathBin)
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

	stdoutStr := strings.TrimSpace(stdout.String())
	if strings.Index(stdoutStr, "failed") != -1 {
		if debug {
			log.Println(stdoutStr)
		}
		err = errors.New(stdoutStr)
		return
	}

	total = len(strings.Split(stdoutStr, "\n"))

	return total, nil
}
