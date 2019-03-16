// Example: collect GPU card metrics and dump into JSON format to stdout.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/OpenCCTV/gpu/gpu"
)

var (
	debug bool
)

func main() {
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}

	metrics, err := metricsGPU.Gets(debug)
	if err != nil {
		panic(err)
	}

	out, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
