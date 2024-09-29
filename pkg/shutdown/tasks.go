package shutdown

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var defaultStopSigs = []os.Signal{syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}

type Tasks struct {
	logger    *slog.Logger
	sigChan   chan os.Signal
	panicChan chan struct{}
	done      chan struct{}
	tasks     []Task
	mu        sync.Mutex
}

type Task func(ctx context.Context) error

func NewShutdownTasks(logger *slog.Logger) (*Tasks, context.Context) {
	appCtx, cancel := context.WithCancel(context.Background())
	t := &Tasks{
		logger:    logger,
		tasks:     make([]Task, 0),
		done:      make(chan struct{}),
		sigChan:   make(chan os.Signal, 1),
		panicChan: make(chan struct{}),
	}
	go func() {
		signal.Notify(t.sigChan, defaultStopSigs...)
		select {
		case sig := <-t.sigChan:
			t.logger.Info(fmt.Sprintf("got stop sig: %v", sig))
		case <-t.panicChan:
		}
		cancel()
		t.gracefulShutdownAll(appCtx)
	}()
	return t, appCtx
}

func (t *Tasks) AddShutdownTask(tasks ...Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasks = append(t.tasks, tasks...)
}

func (t *Tasks) gracefulShutdownAll(ctx context.Context) {
	for i := len(t.tasks) - 1; i >= 0; i-- {
		shutdownFunc := t.tasks[i]
		if shutdownFunc == nil {
			continue
		}
		if err := shutdownFunc(ctx); err != nil {
			t.logger.Info("error while shutting down task", slog.String("error", err.Error()))
		}
		t.tasks[i] = nil
	}
	close(t.done)
}

func (t *Tasks) Wait(panicSource any) {
	if panicSource != nil {
		t.logger.Error(fmt.Sprintf("got panic: %v", panicSource))
		close(t.panicChan)
	}
	<-t.done
	t.logger.Info("gracefully shutdown all tasks")
}

func (t *Tasks) GetSigChan() chan os.Signal {
	return t.sigChan
}
