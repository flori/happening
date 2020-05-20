package happening

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func processorCount() float64 {
	return float64(runtime.NumCPU())
}

func cpuLoad() float64 {
	cpuLoad := 0.0
	cmd := exec.Command("ps", "-A", "-o", "%cpu=0.0")
	cmd.Env = append(os.Environ(), "LANG=C")
	out, err := cmd.Output()
	if err != nil {
		return cpuLoad
	}
	result := strings.Split(string(out), "\n")
	sum := 0.0
	for _, load := range result {
		load := strings.Trim(load, " ")
		if load == "" {
			continue
		}
		f, err := strconv.ParseFloat(load, 64)
		if err != nil {
			continue
		}
		sum += f
	}
	cpuLoad = sum / 100
	return cpuLoad
}

func cpuLoadTotal() float64 {
	return cpuLoad() / processorCount()
}

type LoadTicker struct {
	ticker *time.Ticker
	total  float64
	count  int
	done   chan bool
}

func (lt *LoadTicker) Start() {
	lt.ticker = time.NewTicker(1 * time.Second)
	lt.done = make(chan bool)
	go func() {
		for {
			select {
			case <-lt.done:
				return
			case <-lt.ticker.C:
				lt.total += cpuLoadTotal()
				lt.count++
			}
		}
	}()
}

func normalize(value float64) float32 {
	if value > 1.0 {
		value = 1.0
	}
	return float32(value)
}

func (lt *LoadTicker) Compute() float32 {
	lt.done <- true
	lt.total += cpuLoadTotal()
	lt.count++
	return normalize(lt.total / float64(lt.count))
}
