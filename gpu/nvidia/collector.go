// Collect NVIDIA GPU info via package nvidia-381.
//    add-apt-repository ppa:graphics-drivers/ppa
//    apt-get install nvidia-384
package GPUNvidia

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	//pathBin = "/usr/bin/nvidia-smi"
	pathBin = "nvidia-smi"
)

// 删除输出内容中的\x00和多余的空格
func trimOutput(buffer bytes.Buffer) string {
	return strings.TrimSpace(string(bytes.TrimRight(buffer.Bytes(), "\x00")))
}

func ExecCommand(cmdString string, timeoutInSeconds int) (out []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", cmdString)

	t1 := time.Now() // get current time

	cmd.Env = append(os.Environ(),
		"COLUMNS=512", // overwrite montherfuck monkey hard-coded COLUMNS=80 in golang by default
	)

	elapsed := time.Since(t1)
	log.Printf("bbbbbbb 111 elapsed: %v, [%v]", elapsed, cmdString)

	out, err = cmd.Output()

	elapsed = time.Since(t1)
	log.Printf("bbbbbbb 222 elapsed: %v, [%v], [%v]", elapsed, cmdString, string(out))

	// elapsed := time.Since(t1)
	// log.Printf("bbbbbbb elapsed: %v, [%v]", elapsed, cmdString)

	return
}

func execShell(cmd string, filterFirst bool) (result []string, e error) {
	var timeout = 3
	if len(cmd) == 0 {
		e = fmt.Errorf("cannot run a empty command")
		return
	}
	stdout, _ := ExecCommand(cmd, timeout)

	result = strings.Split(string(bytes.TrimRight(stdout, "\x00")), "\n")
	e = nil

	if filterFirst {
		var i = 0 //过滤掉第一行
		result = append(result[:i], result[i+1:]...)
	}

	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return
}

func Gets(debug bool) (m *[]map[string]interface{}, err error) {
	result := []map[string]interface{}{}
	//total := 8

	total, err := getTotalGPU(debug)
	if err != nil {
		log.Printf("get total gpu error [%v] [%v]", total, err)
		return
	}

	log.Printf("total gpu [%v]", total)

	//	cmd := exec.Command("bash", "-c", "export NV_DRIVER=/var/drivers/nvidia/current;export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$NV_DRIVER/lib:$NV_DRIVER/lib64;export PATH=$PATH:$NV_DRIVER/bin;nvidia-smi -pm 1 > /dev/null && nvidia-smi")
	cmd := exec.Command("bash", "-c", "((nvidia-smi -pm 1 > /dev/null && nvidia-smi) || (export NV_DRIVER=/var/drivers/nvidia/current;export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$NV_DRIVER/lib:$NV_DRIVER/lib64;export PATH=$PATH:$NV_DRIVER/bin;nvidia-smi -pm 1 > /dev/null && nvidia-smi)) | grep -Eiv \"(fail)|(not)\"")

	out, err := cmd.Output()
	if err != nil {
		log.Printf("exce gpu nvidia-smi  fail [%v]", err)
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

/*
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

	arg := fmt.Sprintf("%v --list-gpus", pathBin)
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
					log.Printf("gpu getTotalGPU [%v], [%v]", err, arg)
					return
				}
			}
		} else {
			if debug {
				log.Println(err)
			}
			log.Printf("gpu getTotalGPU [%v]", err)
			return
		}
	}

	if isTimeout {
		errmsg = "run cmd timeout"
		if debug {
			log.Println(errmsg)
		}
		err = errors.New(errmsg)
		log.Printf("gpu getTotalGPU [%v]", err)
		return
	}

	stderrStr := strings.TrimSpace(stderr.String())
	if stderrStr != "" {
		if debug {
			log.Println(stderrStr)
		}
		err = errors.New(stderrStr)
		log.Printf("gpu getTotalGPU [%v]", err)
		return
	}

	stdoutStr := strings.TrimSpace(stdout.String())
	if strings.Index(stdoutStr, "failed") != -1 {
		if debug {
			log.Println(stdoutStr)
		}
		err = errors.New(stdoutStr)
		log.Printf("gpu getTotalGPU [%v]", err)
		return
	}

	total = len(strings.Split(stdoutStr, "\n"))

	return total, nil
}
*/

func getTotalGPU(debug bool) (total int, err error) {
	/*
		var errmsg string
			if _, err = os.Stat(pathBin); os.IsNotExist(err) {
				errmsg = "stat tool not found"
				if debug {
					log.Println(errmsg, pathBin)
				}
				err = errors.New(errmsg)
				return
			}
	*/
	//	arg := fmt.Sprintf("export NV_DRIVER=/var/drivers/nvidia/current;export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$NV_DRIVER/lib:$NV_DRIVER/lib64;export PATH=$PATH:$NV_DRIVER/bin;nvidia-smi --list-gpus")
	arg := fmt.Sprintf("((nvidia-smi --list-gpus) || (export NV_DRIVER=/var/drivers/nvidia/current;export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$NV_DRIVER/lib:$NV_DRIVER/lib64;export PATH=$PATH:$NV_DRIVER/bin;nvidia-smi --list-gpus)) | grep -Eiv \"(fail)|(not)\"")

	strs, err := execShell(arg, false)

	total = len(strs)

	if err != nil {
		return total, err
	}

	return total, nil
}
