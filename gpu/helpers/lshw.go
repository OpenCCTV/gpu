package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	pathBinNvidia = "/usr/bin/nvidia-smi"
	pathBinIntel  = "/usr/bin/intel_gpu_top"
)

// Get display(GPU) venrder info by `lshw`.
func GuessGPUVendors() (vendors []string, err error) {
	//_cmd := `lshw -c display | grep vendor`
	_cmd := `lspci | grep -i 3d`
	out, err := exec.Command("/bin/bash", "-c", _cmd).Output()
	if err != nil {
		log.Println(fmt.Sprintf("[error] issue command failed -%s-", _cmd))

		_, err := os.Stat(pathBinNvidia)
		if err == nil {
			vendors = append(vendors, "nvidia")
		}

		_, err = os.Stat(pathBinIntel)
		if err == nil {
			vendors = append(vendors, "intel")
		}

		if len(vendors) > 0 {
			return vendors, nil
		}

		err = errors.New("issue lshw failed and both nvidia and intel collector tool not found")
		return vendors, err

	} else {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			columns := strings.Split(line, ":")
			if len(columns) != 2 {
				log.Println("parse vendor info failed", line, columns)
				continue
			}

			vendor := strings.TrimSpace(columns[1])
			vendor = strings.ToLower(vendor)
			if strings.Index(vendor, "intel") != -1 {
				vendor = "intel"
			} else if strings.Index(vendor, "nvidia") != -1 {
				vendor = "nvidia"
			}

			if vendor != "" {
				vendors = append(vendors, vendor)
			}
		}
	}

	return vendors, nil
}
