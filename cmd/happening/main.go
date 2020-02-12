package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	happening "github.com/flori/happening"
)

var config happening.Config

func init() {
	d := happening.NewConfig()
	flag.StringVar(&config.Name, "name", d.Name, "name of the event that happened")
	flag.StringVar(&config.ReportURL, "report", d.ReportURL, "send events to this report URL")
	flag.BoolVar(&config.StoreReport, "store-report", d.StoreReport, "store the report iff true")
	flag.StringVar(&config.SuccessCode, "success", d.SuccessCode, "consider these exit codes (separated by ,) as success")
	flag.StringVar(&config.PingURL, "ping", d.PingURL, "ping URL after successful execution of command")
	flag.StringVar(&config.FlagHostname, "hostname", d.FlagHostname, "overwrite os hostname with this value")
	flag.UintVar(&config.Retries, "retries", d.Retries, "retry requests that many times")
	flag.DurationVar(&config.RetryDelay, "retry-delay", d.RetryDelay, "delay for this duration between retries")
	flag.BoolVar(&config.CollectOutput, "collect-output", d.CollectOutput, "collect output of executed command")
	flag.BoolVar(&config.SuppressOutput, "suppress-output", d.SuppressOutput, "suppress output of executed command")
	flag.StringVar(&config.Chdir, "cd", d.Chdir, "change directory to here before running command")
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
