package helpers

import (
	"errors"
	"log"
	"os/exec"
	"strings"
)

// Get display(GPU) venrder info by `lshw`.
func GuessGPUVendors() (vendors []string, err error) {
	var errmsg string
	out, err := exec.Command("bash", "-c", `lshw -c display | grep vendor`).Output()
	if err != nil {
		errmsg = "get lshw info failed"
		log.Println(errmsg, err)
		return vendors, errors.New(errmsg)
	}

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

	return vendors, nil
}
