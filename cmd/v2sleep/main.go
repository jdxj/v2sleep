package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/config"
	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/router"
)

var (
	confPath = flag.String("conf", "", "conf path")
)

func main() {
	flag.Parse()

	if *confPath == "" {
		logrus.Fatalf("conf not set")
	}

	conf, err := config.ReadConfig(*confPath)
	if err != nil {
		logrus.Fatalf("read config file err: %s", err)
	}

	dao.Init(conf.DB)
	router.Run(conf.Web)
}
