// Example: collect GPU card metrics and dump into JSON format to stdout.
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/MonitorMetrics/gpu/gpu"
)

var (
	Debug = flag.Bool("debug", false, "")
)

func main() {
	flag.Parse()
	metrics, err := metricsGPU.Gets(*Debug)
	if err != nil {
		panic(err)
	}

	out, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
