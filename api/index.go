package api

import (
	"gitlab.com/hieuhani/permitbox/app"
	"gitlab.com/hieuhani/permitbox/pkg/shutdown"
	"log/slog"
	"net/http"
	"os"
)

var (
	application *app.Application
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tasks, _ := shutdown.NewShutdownTasks(logger)
	var err error
	application, err = app.NewApplication(logger, tasks)
	if err != nil {
		logger.Error("cannot create application", err)
		os.Exit(1)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	application.GetHttpHandler().ServeHTTP(w, r)
}
