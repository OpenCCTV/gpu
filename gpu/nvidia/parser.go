// Parse `nvidia-384` `nvidia-smi` program output into map.
package GPUNvidia

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var percentRe = regexp.MustCompile(`([0-9]+)%`)
var tempRe = regexp.MustCompile(`([0-9]+)C`)
var sizeRe = regexp.MustCompile(`([0-9]+)MiB`)

func parseFan(result *map[string]interface{}, line string) {
	columns := strings.Fields(line)
	(*result)["fan"] = -1

	if len(columns) <= 0 {
		log.Println("parse columns failed")
		return
	}

	if columns[0] == "N/A" {
		// some machine return N/A
		return
	}

	subMatch := percentRe.FindStringSubmatch(columns[0])
	if len(subMatch) != 2 {
		log.Println(fmt.Sprintf("match re failed -%s-", columns[0]))
		return
	}

	percent, err := strconv.Atoi(subMatch[1])
	if err != nil {
		log.Println("parse into int failed", subMatch, err)
		return
	}

	(*result)["fan"] = percent
}

func parseTemp(result *map[string]interface{}, line string) {
	columns := strings.Fields(line)
	if len(columns) < 2 {
		log.Println("parse columns failed")
		return
	}
	subMatch := tempRe.FindStringSubmatch(columns[1])
	if len(subMatch) != 2 {
		log.Println("match re failed", tempRe, columns[1])
		return
	}

	temp, err := strconv.Atoi(subMatch[1])
	if err != nil {
		log.Println("parse into int failed", subMatch, err)
		return
	}

	(*result)["temp"] = temp
}

func parseMemoryUsage(result *map[string]interface{}, line string) {
	columns := strings.Split(line, "/")
	if len(columns) != 2 {
		log.Println("parse columns failed")
		return
	}

	var subMatch []string
	subMatch = sizeRe.FindStringSubmatch(strings.TrimSpace(columns[0]))
	if len(subMatch) != 2 {
		log.Println("match re failed", sizeRe, columns[0])
		return
	}

	used, err := strconv.Atoi(subMatch[1])
	if err != nil {
		log.Println("parse into int failed", subMatch[0], err)
		return
	}

	subMatch = sizeRe.FindStringSubmatch(strings.TrimSpace(columns[1]))
	if len(subMatch) != 2 {
		log.Println("match re failed", sizeRe, columns[1])
		return
	}

	total, err := strconv.Atoi(subMatch[1])
	if err != nil {
		log.Println("parse into int failed", subMatch[1], err)
		return
	}

	var usedPercent float64
	if total > 0 {
		usedPercent = float64(used) / float64(total) * 100
	} else {
		usedPercent = 0
	}

	(*result)["memory_usage_percent"] = usedPercent

	// convert MiB into MB, 1MiB = 1048576 bytes
	// https://en.wikipedia.org/wiki/Mebibyte
	MiBInBytes := 1048576
	(*result)["memory_usage"] = used * MiBInBytes
	(*result)["memory_total"] = total * MiBInBytes

}

func parseGPUUtil(result *map[string]interface{}, line string) {
	columns := strings.Fields(line)
	if len(columns) < 2 {
		log.Println("parse columns failed")
		return
	}

	subMatch := percentRe.FindStringSubmatch(columns[0])
	if len(subMatch) != 2 {
		log.Println("match re failed", percentRe, columns[0])
		return
	}

	percent, err := strconv.Atoi(subMatch[1])
	if err != nil {
		log.Println("parse into int failed", subMatch, err)
		return
	}

	(*result)["util"] = percent
}

func ParseGPURow(line string) (*map[string]interface{}, error) {
	result := map[string]interface{}{}

	columns := strings.Split(line, "|")

	columnsFiltered := []string{}
	for _, column := range columns {
		column = strings.TrimSpace(column)
		if column != "" {
			columnsFiltered = append(columnsFiltered, column)
		}
	}
	if len(columnsFiltered) < 3 {
		return nil, errors.New("parse failed")
	} else {
		parseFan(&result, columnsFiltered[0])
		parseTemp(&result, columnsFiltered[0])
		parseMemoryUsage(&result, columnsFiltered[1])
		parseGPUUtil(&result, columnsFiltered[2])
	}

	return &result, nil
}
