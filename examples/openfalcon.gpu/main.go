// Example Collect GPU cards metrics and post to open-falcon agent HTTP PUSH API.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MonitorMetrics/falcon_helpers/agent"
	"github.com/MonitorMetrics/falcon_helpers/model"
	"github.com/MonitorMetrics/gpu/gpu"
)

var (
	Debug           = flag.Bool("debug", false, `enable debug`)
	FalconURL       = flag.String("url", "http://127.0.0.1:1988/v1/push", "falcon agnet PUSH URL")
	IntervalCollect = flag.Int("interval", 60, "collect stat interval in seconds")
	ShortenInterval = flag.Bool("s", false, "shorten interval to 5 seconds for debug")
)

func funcCallback(m *[]map[string]interface{}, args ...interface{}) {
	metrics := []*modelMetric.MetricItem{}

	now := time.Now().Unix()

	for _, item := range *m {
		for k, v := range item {
			if k == "vendor" || k == "gpuid" {
				continue
			}
			tags := fmt.Sprintf("vendor=%v,gpuid=%v", item["vendor"], item["gpuid"])
			item := modelMetric.MetricItem{
				Endpoint:    "127.0.0.1",
				Metric:      metricsGPU.MetricPrefix + k,
				Value:       v,
				CounterType: "GAUGE",
				Tags:        tags,
				Timestamp:   now,
				Step:        60,
			}
			metrics = append(metrics, &item)
		}
	}

	out, err := json.Marshal(metrics)
	if err != nil {
		log.Println(err)
		return
	}

	respBody := helperAgent.SendToFalconAgent(*FalconURL, string(out))
	if *Debug {
		log.Println("falcon respBody", respBody)
	}
}

func main() {
	flag.Parse()

	if *Debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}

	intervalInSeconds := *IntervalCollect
	if *ShortenInterval {
		intervalInSeconds = 5
	}

	for {
		metrics, err := metricsGPU.Gets(*Debug)
		if err != nil {
			log.Println(err)
		} else {

			go funcCallback(metrics)
		}
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
	}

}
