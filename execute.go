package happening

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var load = LoadTicker{}

func getSuccessCodes(config *Config) []int {
	var codes []int
	if config.SuccessCode == "" {
		return codes
	}
	for _, code := range strings.Split(config.SuccessCode, ",") {
		c, err := strconv.ParseInt(code, 10, 32)
		if err != nil {
			log.Fatalf("invalid exit code, %v", err)
		}
		codes = append(codes, int(c))
	}
	return codes
}

func isSuccess(exitCode int, config *Config) bool {
	codes := getSuccessCodes(config)
	for _, code := range codes {
		if code == exitCode {
			return true
		}
	}
	return false
}

func determineHostname(flagHostname string) string {
	if flagHostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		return hostname
	}
	return flagHostname
}

func determineUser(configUser string) string {
	if configUser == "" {
		configUser = os.Getenv("LOGNAME")
	}
	return configUser
}

func setEventFields(config Config, event *Event) {
	if config.Started != "" {
		started, err := time.Parse(time.RFC3339, config.Started)
		if err != nil {
			log.Fatalf("invalid time format (should be RFC3339), %v", err)
		}
		event.Started = started
	}
	if config.Duration != time.Duration(0) {
		event.Duration = config.Duration
		if config.Started == "" {
			event.Started = event.Started.Add(-event.Duration)
		}
	}
	if config.Output != "" {
		event.Output = config.Output
	}
}

func Execute(config Config, block func(output io.Writer) bool) *Event {
	hostname := determineHostname(config.Hostname)
	user := determineUser(config.User)
	started := time.Now()
	load.Start()
	success := true
	duration := time.Duration(0)
	outputString := ""
	pid := os.Getpid()

	pi := getProcinfoSelf()

	if block != nil {
		if config.CollectOutput {
			output := new(bytes.Buffer)
			success = block(output)
			duration = time.Now().Sub(started)
			outputString = output.String()
		} else {
			success = block(NullWriter)
			duration = time.Now().Sub(started)
		}
	}

	event := Event{
		Id:          GenerateUUIDv4(),
		Name:        config.Name,
		Context:     config.Context,
		Output:      outputString,
		Started:     started,
		Duration:    duration,
		Success:     success,
		Hostname:    hostname,
		User:        user,
		Pid:         pid,
		Load:        load.Compute(),
		CpuUsage:    pi.CpuUsage,
		MemoryUsage: pi.MemoryUsage,
		Store:       config.StoreReport,
	}

	setEventFields(config, &event)

	return &event
}

func ExecuteCmd(config Config, cmdArgs []string) *Event {
	hostname := determineHostname(config.Hostname)
	user := determineUser(config.User)
	if len(cmdArgs) > 0 {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		var outputBuffer bytes.Buffer
		stat, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("error while stating stdin: %v", err)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			cmd.Stdin = os.Stdin
		}
		if config.CollectOutput {
			if config.SuppressOutput {
				cmd.Stdout = io.Writer(&outputBuffer)
				cmd.Stderr = io.Writer(&outputBuffer)
			} else {
				mwriter := io.MultiWriter(os.Stdout, &outputBuffer)
				cmd.Stdout = mwriter
				cmd.Stderr = mwriter
			}
		} else {
			if config.SuppressOutput {
				cmd.Stdout = ioutil.Discard
				cmd.Stderr = ioutil.Discard
			} else {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stdout
			}
		}
		var oldDir string
		var success bool
		var exitCode int
		var signal string
		var output string
		var pi *procinfo
		pid := 0
		if config.Chdir != "" {
			oldDir, err = os.Getwd()
			if err == nil {
				err = os.Chdir(config.Chdir)
			}
		}
		started := time.Now()
		load.Start()
		if err == nil {
			err = cmd.Start()
		}
		if err == nil {
			pid = cmd.Process.Pid
			err = cmd.Wait()
			processState := cmd.ProcessState
			pi = getProcinfo(processState)
			exitCode = processState.ExitCode()
			status := processState.Sys().(syscall.WaitStatus)
			if status.Signaled() {
				signal = status.Signal().String()
			}
			success = isSuccess(exitCode, &config)
			output = outputBuffer.String()
		} else {
			output = fmt.Sprintf(
				"happening: Starting \"%s\" failed with error \"%v\"", cmdArgs[0], err)
			log.Println(output)
			exitCode = 1
			success = false
		}
		if config.Chdir != "" {
			os.Chdir(oldDir)
		}
		event := Event{
			Id:          GenerateUUIDv4(),
			Name:        config.Name,
			Context:     config.Context,
			Command:     cmdArgs,
			Output:      output,
			Started:     started,
			Duration:    time.Now().Sub(started),
			Success:     success,
			ExitCode:    exitCode,
			Signal:      signal,
			Hostname:    hostname,
			User:        user,
			Pid:         pid,
			Load:        load.Compute(),
			CpuUsage:    pi.CpuUsage,
			MemoryUsage: pi.MemoryUsage,
			Store:       config.StoreReport,
		}
		setEventFields(config, &event)
		return &event
	} else {
		return Execute(config, nil)
	}
}
