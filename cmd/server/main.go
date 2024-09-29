package main

import (
	"gitlab.com/hieuhani/permitbox/app"
	"gitlab.com/hieuhani/permitbox/pkg/shutdown"
	"log/slog"
	"os"
	"runtime/debug"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tasks, _ := shutdown.NewShutdownTasks(logger)
	application, err := app.NewApplication(logger, tasks)
	handleError(logger, err)
	err = application.ServeHTTP()
	handleError(logger, err)
}

func handleError(logger *slog.Logger, err error) {
	if err != nil {
		trace := debug.Stack()
		logger.Error("cannot start application", slog.String("error", err.Error()), slog.String("stack", string(trace)))
		os.Exit(1)
	}
}
