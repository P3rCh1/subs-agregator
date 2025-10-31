package subs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *ServerAPI) Create(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (s *ServerAPI) Read(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (s *ServerAPI) Update(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (s *ServerAPI) Delete(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (s *ServerAPI) List(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}
