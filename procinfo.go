package happening

import (
	"os"
	"runtime"
	"syscall"
	"time"
)

type procinfo struct {
	CpuUsage    float64
	MemoryUsage float64
}

func maxRSSMultiplier() float64 {
	switch runtime.GOOS {
	case "linux":
		return 1 << 10
	default:
		return 1
	}
}

func computeCPUTime(rusage *syscall.Rusage) float64 {
	userCpuTime := time.Duration(rusage.Utime.Sec)*time.Second +
		time.Duration(rusage.Utime.Usec)*time.Microsecond
	sysCpuTime := time.Duration(rusage.Stime.Sec)*time.Second +
		time.Duration(rusage.Stime.Usec)*time.Microsecond
	cpuTime := userCpuTime + sysCpuTime
	return cpuTime.Seconds()
}

func getProcinfo(processState *os.ProcessState) *procinfo {
	rusage := processState.SysUsage().(*syscall.Rusage)
	return &procinfo{
		CpuUsage:    computeCPUTime(rusage),
		MemoryUsage: float64(rusage.Maxrss),
	}
}

func getProcinfoSelf() *procinfo {
	rusage := new(syscall.Rusage)
	err := syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
	if err == nil {
		return &procinfo{
			CpuUsage:    computeCPUTime(rusage),
			MemoryUsage: float64(rusage.Maxrss) * maxRSSMultiplier(),
		}
	}
	return &procinfo{}
}
