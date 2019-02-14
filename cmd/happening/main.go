package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	happening "github.com/flori/happening"
)

var config happening.Config

func init() {
	flag.StringVar(&config.Name, "name", "some event", "name of the event that happened")
	flag.StringVar(&config.ReportURL, "report", "", "send events to this report URL")
	flag.BoolVar(&config.StoreReport, "store-report", true, "store the report iff true")
	flag.StringVar(&config.SuccessCode, "success", "0", "consider these exit codes (separated by ,) as success")
	flag.StringVar(&config.PingURL, "ping", "", "ping URL after successful execution of command")
	flag.StringVar(&config.FlagHostname, "hostname", "", "overwrite os hostname with this value")
	flag.UintVar(&config.Retries, "retries", 3, "retry requests that many times")
	flag.DurationVar(&config.RetryDelay, "retry-delay", time.Second, "delay for this duration between retries")
	flag.BoolVar(&config.CollectOutput, "collect-output", false, "collect output of executed command")
	flag.BoolVar(&config.SuppressOutput, "suppress-output", false, "suppress output of executed command")
	flag.StringVar(&config.Chdir, "cd", "", "change directory to here before running command")
	quiet := flag.Bool("quiet", false, "don't output log messages in success case")
	flag.Parse()
	if *quiet {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	cmdArgs := flag.Args()
	event := happening.ExecuteCmd(config, cmdArgs)
	if event.Success && config.PingURL != "" {
		happening.Ping(&config)
	}
	if config.ReportURL != "" {
		happening.SendEvent(event, &config)
	}
	log.Printf("\"%s\" took %s: %v",
		event.Name, event.Duration, string(happening.EventToJSON(event)))
	os.Exit(event.ExitCode)
}
