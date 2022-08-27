package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jdxj/v2sleep/model"
)

func GetShare(ctx *gin.Context) {
	sig := strings.TrimSpace(ctx.Query("sig"))
	if sig == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sigRsp, err := model.GetSig(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if sig != sigRsp.Sig {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	buf, err := model.GenShare(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(buf.Len()), "text/html", buf, nil)
}
