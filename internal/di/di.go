package di

import (
	"github.com/samber/do"
	namespaceHttp "gitlab.com/hieuhani/permitbox/internal/core/namespace/transport/http"
	"gitlab.com/hieuhani/permitbox/internal/handler"
	"log/slog"
)

func NewInjector(logger *slog.Logger) *do.Injector {
	injector := do.New()
	do.ProvideValue(injector, logger)

	do.Provide(injector, NewBaseHandler)

	do.Provide(injector, NewNamespaceHandler)

	return injector
}

func NewBaseHandler(i *do.Injector) (handler.BaseHandler, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	return handler.NewBaseHandler(logger), nil
}

func NewNamespaceHandler(injector *do.Injector) (namespaceHttp.NamespaceHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](injector)
	return namespaceHttp.NewNamespaceHandler(baseHandler), nil
}
