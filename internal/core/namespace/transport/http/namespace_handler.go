package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/hieuhani/permitbox/internal/handler"
	"net/http"
)

type NamespaceHandler struct {
	handler.BaseHandler
}

func NewNamespaceHandler(baseHandler handler.BaseHandler) NamespaceHandler {
	return NamespaceHandler{
		BaseHandler: baseHandler,
	}
}

func (h NamespaceHandler) GetAllNamespaces(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{
		Data: "hello world",
	})
}
