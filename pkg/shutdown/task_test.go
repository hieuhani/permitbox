package shutdown

import (
	"bytes"
	"context"
	"log/slog"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewShutdownTasks(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 0)), nil))
	tasks, _ := NewShutdownTasks(logger)
	assert.NotNil(t, tasks)
	val := 0
	tasks.AddShutdownTask(
		func(ctx context.Context) error {
			val--
			return nil
		},
		func(ctx context.Context) error {
			val *= 2
			return nil
		},
		func(ctx context.Context) error {
			val += 2
			return nil
		},
	)
	go func() {
		tasks.GetSigChan() <- syscall.SIGINT
	}()
	waitDone := make(chan struct{})
	go func() {
		tasks.Wait(nil)
		close(waitDone)
	}()
	select {
	case <-waitDone:
	case <-time.After(100 * time.Millisecond):
	}
	assert.Equal(t, 3, val)
}
