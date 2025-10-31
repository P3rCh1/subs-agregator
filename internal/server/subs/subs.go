package subs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*APIServer) Create(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (*APIServer) Read(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (*APIServer) Update(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (*APIServer) Delete(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}

func (*APIServer) List(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented\n")
}
