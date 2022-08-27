package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	User string
	Pass string
	Port int
}

func Run(conf Config) {
	r := gin.Default()
	register(conf, r)

	err := r.Run(fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		logrus.Errorf("run: %s", err)
	}
}

func getBasicAuth(conf Config) gin.HandlerFunc {
	conf.User = strings.TrimSpace(conf.User)
	conf.Pass = strings.TrimSpace(conf.Pass)
	if conf.User == "" || conf.Pass == "" {
		logrus.Fatalf("account not set")
	}

	return gin.BasicAuth(gin.Accounts{
		conf.User: conf.Pass,
	})
}

func register(conf Config, root gin.IRouter) {
	root.GET("/", hello)
	rShare := root.Group("/share")
	{
		rShare.GET("/", GetShare)
	}

	rSections := root.Group("/sections", getBasicAuth(conf))
	{
		rSection := rSections.Group("/section")
		{
			rSection.POST("/", SetSig)
			rSection.GET("/", GetSig)
		}
	}

	rSubConfig := root.Group("/sub_config", getBasicAuth(conf))
	{
		rSubConfig.GET("/", ListSubConfig)
		rSubConfig.DELETE("/", DeleteSubConfig)
		// rSubConfig.POST("/clash_sub_addr", AddClashSubAddr)
		rSubConfig.POST("/v2ray_sub_addr", AddV2raySubAddr)
	}
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world!")
}
