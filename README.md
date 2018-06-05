# GPU card metrics for monitoring.

We do not parse values in `int` or `float*` type and keep them original, end-user should take care type by itself.
For example, Open-falcon will parses all values into JSON `float64` finally.


## Tags

| KEY | TYPE | VALUE | NOTES |
|-----|------|-------|-------|
| gpuid | string | 0~n | Intel GPU card is hard-coded 0 |
| vendor | string | "intel" "nvidia" |

## Metrics

Intel GPU Card Metrics

| KEY | TYPE | VALUE | NOTES |
|-----|------|-------|-------|
| renderbusy | GAUGE | 0~100 | in percent |


NVIDIA GPU Card Metrics

| KEY | TYPE | VALUE | NOTES |
|-----|------|-------|-------|
| fan | GAUGE | int  |
| temp | GAUGE | int  |
| util | GAUGE | 0~100  | in percent |
| memory_usage | GAUGE | int | in bytes |
| memory_total | GAUGE | int | in bytes |
| memory_usage_percent | GAUGE  | memory_usage/memory_total*100 | in percent |

## Build

Collect Intel GPU cards stats requried intel-gpu-tools `/usr/bin/intel_gpu_top`
 
    sudo apt-get install intel-gpu-tools


NOTICE: Intel GPU card stats collecotr required root privilege.


Collect NVIDIA GPU cards stats requried nvidia-352 `/usr/bin/nvidia-smi`
 
    sudo apt-get install nvidia-352


Build example

	go get -v github.com/MonitorMetrics/redis
	cd $GOPATH/src/github.com/MonitorMetrics/redis/examples/json.gpu
	go build -o json.gpu.bin


Start in debug mode

	sudo ./json.gpu.bin -debug

For more detail, see source code.
