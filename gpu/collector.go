package metricsGPU

import (
	"github.com/MonitorMetrics/gpu/gpu/helpers"
	"github.com/MonitorMetrics/gpu/gpu/intel"
	"github.com/MonitorMetrics/gpu/gpu/nvidia"
	"log"
)

type FuncGets func(debug bool) (metrics *[]map[string]interface{}, err error)

var (
	funcMap = map[string]FuncGets{
		"intel":  GPUIntel.Gets,
		"nvidia": GPUNvidia.Gets,
	}
)

func Gets(debug bool) (*[]map[string]interface{}, error) {
	vendors, err := helpers.GuessGPUVendors()
	if err != nil {
		if debug {
			log.Fatalln("GuessGPUVendors failed", err)
		}
		return nil, err
	}

	uniq := map[string]bool{}

	merge := []map[string]interface{}{}

	for _, vendor := range vendors {
		if _, ok := uniq[vendor]; ok {
			continue
		}

		f, ok := funcMap[vendor]
		if !ok {
			if debug {
				log.Println("vendor not support", vendor)
			}
			continue
		}

		metrics, err := f(debug)
		if err != nil {
			if debug {
				log.Println("Gets failed", err)
			}
		} else {
			for _, metric := range *metrics {
				merge = append(merge, metric)
			}
		}
	}

	return &merge, nil
}
