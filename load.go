package happening

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func load() float32 {
	return float32(loadAvg() / processorCount())
}
