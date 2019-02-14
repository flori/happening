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

func exitCode(err error) int {
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
	}
	return 0
}

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

func Execute(config Config, block func() bool) *Event {
	hostname := determineHostname(config.FlagHostname)
	started := time.Now()
	success := true
	duration := time.Duration(0)
	if block != nil {
		success = block()
		duration = time.Now().Sub(started)
	}
	return &Event{
		Id:       GenerateUUIDv4(),
		Name:     config.Name,
		Started:  started,
		Duration: duration,
		Success:  success,
		Hostname: hostname,
		Pid:      os.Getpid(),
		Store:    config.StoreReport,
	}
}

func ExecuteCmd(config Config, cmdArgs []string) *Event {
	hostname := determineHostname(config.FlagHostname)
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
		var ec int
		var output string
		pid := 0
		if config.Chdir != "" {
			oldDir, err = os.Getwd()
			if err == nil {
				err = os.Chdir(config.Chdir)
			}
		}
		started := time.Now()
		if err == nil {
			err = cmd.Start()
		}
		if err == nil {
			err = cmd.Wait()
			ec = exitCode(err)
			success = isSuccess(ec, &config)
			output = outputBuffer.String()
			pid = cmd.ProcessState.Pid()
		} else {
			output = fmt.Sprintf(
				"happening: Starting \"%s\" failed with error \"%v\"", cmdArgs[0], err)
			log.Println(output)
			ec = 1
			success = false
		}
		if config.Chdir != "" {
			os.Chdir(oldDir)
		}
		return &Event{
			Id:       GenerateUUIDv4(),
			Name:     config.Name,
			Command:  cmdArgs,
			Output:   output,
			Started:  started,
			Duration: time.Now().Sub(started),
			Success:  success,
			ExitCode: ec,
			Hostname: hostname,
			Pid:      pid,
			Store:    config.StoreReport,
		}
	} else {
		return Execute(config, nil)
	}
}
