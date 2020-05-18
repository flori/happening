package happening

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func processorCount() float64 {
	processorCount := 0
	out, err := exec.Command("getconf", "_NPROCESSORS_ONLN").Output()
	if err != nil {
		return 0.0

	}
	i, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 32)
	if err != nil {
		return 0.0
	}
	processorCount = int(i)
	return float64(processorCount)
}

func loadAvg() float64 {
	loadAvg := 0.0
	cmd := exec.Command("uptime")
	cmd.Env = append(os.Environ(), "LANG=C")
	out, err := cmd.Output()
	if err != nil {
		return 0.0
	}
	result := strings.Split(string(out), " ")
	l := len(result)
	if l > 3 {
		loadValue := strings.Trim(result[l-3], ",")
		f, err := strconv.ParseFloat(loadValue, 64)
		if err != nil {
			return 0.0
		}
		loadAvg = f
	} else {
		return 0.0
	}
	return loadAvg
}

func currentLoadAverage() float64 {
	return loadAvg() / processorCount()
}

type LoadTicker struct {
	ticker *time.Ticker
	total  float64
	count  int
	done   chan bool
}

func (lt *LoadTicker) Start() {
	lt.ticker = time.NewTicker(1 * time.Minute)
	lt.done = make(chan bool)
	go func() {
		for {
			select {
			case <-lt.done:
				return
			case <-lt.ticker.C:
				lt.total += currentLoadAverage()
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
	lt.total += currentLoadAverage()
	lt.count++
	return normalize(lt.total / float64(lt.count))
}
