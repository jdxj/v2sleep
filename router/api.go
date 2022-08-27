package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type handler func(ctx context.Context) (interface{}, error)

func handle(ctx *gin.Context, req interface{}, h handler) {
	if req != nil {
		err := ctx.ShouldBindJSON(req)
		if err != nil {
			reply(ctx, http.StatusBadRequest, nil, err)
			return
		}
	}

	c, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rsp, err := h(c)
	// todo: 自定义错误
	reply(ctx, 0, rsp, err)
}

type response struct {
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func reply(ctx *gin.Context, code int, data interface{}, err error) {
	rsp := response{
		Data: data,
	}
	if err != nil {
		rsp.Msg = err.Error()
	}
	if code == 0 && err != nil {
		code = http.StatusInternalServerError
	}
	ctx.JSON(code, rsp)
}
