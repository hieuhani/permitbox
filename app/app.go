package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"gitlab.com/hieuhani/permitbox/asset"
	"gitlab.com/hieuhani/permitbox/internal/config"
	namespaceHttp "gitlab.com/hieuhani/permitbox/internal/core/namespace/transport/http"
	"gitlab.com/hieuhani/permitbox/internal/di"
	"gitlab.com/hieuhani/permitbox/pkg/database"
	"gitlab.com/hieuhani/permitbox/pkg/shutdown"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 1 * time.Second
)

type Application struct {
	Logger   *slog.Logger
	Tasks    *shutdown.Tasks
	Injector *do.Injector
	Config   config.AppConfig
}

func NewApplication(logger *slog.Logger, tasks *shutdown.Tasks) (*Application, error) {
	cfg, err := config.InitConfig[config.AppConfig](asset.EmbeddedFiles)
	if err != nil {
		return nil, err
	}
	getDbFunc, atomicExecutor, err := database.New(cfg.Db, tasks, asset.EmbeddedFiles)
	if err != nil {
		return nil, err
	}
	injector := di.NewInjector(logger)
	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, tasks)
	do.ProvideValue(injector, getDbFunc)
	do.ProvideValue(injector, atomicExecutor)
	return &Application{
		Logger:   logger,
		Tasks:    tasks,
		Injector: injector,
		Config:   cfg,
	}, nil
}

func (a *Application) GetHttpHandler() http.Handler {
	namespaceHandler := do.MustInvoke[namespaceHttp.NamespaceHandler](a.Injector)

	r := gin.New()
	v1Routes := r.Group("/api/v1")

	namespaceRoutes := v1Routes.Group("/namespaces")
	namespaceRoutes.GET("", namespaceHandler.GetAllNamespaces)
	return r
}

func (a *Application) ServeHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.HttpPort),
		Handler:      a.GetHttpHandler(),
		ErrorLog:     log.New(os.Stderr, "", 0),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	a.Tasks.AddShutdownTask(
		func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, defaultShutdownPeriod)
			defer cancel()
			return srv.Shutdown(ctx)
		},
	)

	a.Logger.Info(fmt.Sprintf("starting server on %s", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	a.Logger.Info(fmt.Sprintf("stopped server on %s", srv.Addr))

	return nil
}
