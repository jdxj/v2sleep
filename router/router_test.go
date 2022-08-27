package router

import (
	"os"
	"testing"

	"github.com/jdxj/v2sleep/dao"
)

var (
	c = Config{
		User: "jdxj",
		Pass: "123456",
		Port: 8080,
	}

	dbC = dao.Config{
		Host:   "192.168.50.200",
		Port:   3306,
		User:   "root",
		Pass:   "123456",
		DBName: "v2sleep",
	}
)

func TestMain(t *testing.M) {
	dao.Init(dbC)
	os.Exit(t.Run())
}

func TestRun(t *testing.T) {
	Run(c)
}
