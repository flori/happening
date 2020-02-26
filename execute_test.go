package happening

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecuteNil(t *testing.T) {
	result := Execute(Config{Name: "test"}, nil)
	assert.Equal(t, "test", result.Name, "correct name")
	assert.True(t, result.Success, "successful if nil")
	assert.Equal(t, time.Duration(0), result.Duration, "0 duration")
	assert.Equal(t, 0, result.ExitCode, "exit code not set")
}

func TestExecuteBlock(t *testing.T) {
	result := Execute(Config{
		Name:          "test",
		CollectOutput: true,
	}, func(output io.Writer) bool {
		time.Sleep(time.Second)
		fmt.Fprintln(output, "hello world")
		return false
	})
	assert.Equal(t, "test", result.Name, "correct name")
	assert.False(t, result.Success, "not successful if block returns false")
	assert.Condition(
		t,
		func() bool { return result.Duration >= time.Second },
		">= 1s duration",
	)
	assert.Equal(t, 0, result.ExitCode, "exit code not set")
	assert.Equal(t, "hello world\n", result.Output, "collected output")
}

func TestExecuteBlockSuppressOutput(t *testing.T) {
	result := Execute(Config{Name: "test"}, func(output io.Writer) bool {
		time.Sleep(time.Second)
		fmt.Fprintln(output, "hello world")
		return false
	})
	assert.Equal(t, "test", result.Name, "correct name")
	assert.False(t, result.Success, "not successful if block returns false")
	assert.Equal(t, 0, result.ExitCode, "exit code not set")
	assert.Empty(t, result.Output, "collected output")
}

func TestExecuteCmdFail(t *testing.T) {
	result := ExecuteCmd(
		Config{Name: "test"},
		[]string{"sh", "-c", "sleep 1; exit 23"},
	)
	assert.Equal(t, "test", result.Name, "correct name")
	assert.False(t, result.Success, "not successful if cmd returns != 0")
	assert.Condition(
		t,
		func() bool { return result.Duration >= time.Second },
		">= 1s duration",
	)
	assert.Equal(t, 23, result.ExitCode, "exit code set")
}

func TestExecuteCmdSuccess(t *testing.T) {
	result := ExecuteCmd(
		Config{Name: "test", SuccessCode: "42,23", CollectOutput: true},
		[]string{"sh", "-c", `echo "hello world"; exit 23`},
	)
	assert.Equal(t, "test", result.Name, "correct name")
	assert.True(t, result.Success, "not successful if cmd returns != 0")
	assert.Equal(t, 23, result.ExitCode, "exit code set")
	assert.Equal(t, "hello world\n", result.Output, "collected output")
}

func TestExecuteCmdSuccessSuppressOutput(t *testing.T) {
	result := ExecuteCmd(
		Config{
			Name:           "test",
			SuccessCode:    "0",
			CollectOutput:  true,
			SuppressOutput: true,
		},
		[]string{"sh", "-c", `echo "hello world"`},
	)
	assert.Equal(t, "test", result.Name, "correct name")
	assert.True(t, result.Success, "not successful if cmd returns != 0")
	assert.Equal(t, "hello world\n", result.Output, "collected output")
}

func TestExecuteCmdSuccessSuppressOutput2(t *testing.T) {
	result := ExecuteCmd(
		Config{
			Name:           "test",
			SuccessCode:    "0",
			CollectOutput:  false,
			SuppressOutput: true,
		},
		[]string{"sh", "-c", `echo "hello world"`},
	)
	assert.Equal(t, "test", result.Name, "correct name")
	assert.True(t, result.Success, "not successful if cmd returns != 0")
	assert.Empty(t, result.Output, "no collected output")
}

func TestSetEventFieldsStarted(t *testing.T) {
	event := Event{}
	setEventFields(
		Config{},
		&event,
	)
	assert.Equal(t, time.Time{}, event.Started, "is zero time unless Started was set")
	timeString := "2020-02-26T15:27:10+01:00"
	setEventFields(
		Config{
			Started: timeString,
		},
		&event,
	)
	time, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, time, event.Started, "is configured time value")
}

func TestSetEventFieldsDuration(t *testing.T) {
	event := Event{}
	setEventFields(
		Config{},
		&event,
	)
	assert.Equal(t, time.Duration(0), event.Duration, "is zero duration unless Duration was set")
	duration, err := time.ParseDuration("6m66s")
	if err != nil {
		t.Error(err)
	}
	setEventFields(
		Config{
			Duration: duration,
		},
		&event,
	)
	assert.Equal(t, duration, event.Duration, "is configured duration value")
}

func TestSetEventFieldsOutput(t *testing.T) {
	event := Event{}
	setEventFields(
		Config{},
		&event,
	)
	assert.Equal(t, "", event.Output, "is zero duration unless Output was set")
	output := "hello world"
	setEventFields(
		Config{
			Output: output,
		},
		&event,
	)
	assert.Equal(t, output, event.Output, "is configured output value")
}
