package dao

import (
	"os"
	"testing"
)

var (
	c = Config{
		Host:   "192.168.50.200",
		Port:   3306,
		User:   "root",
		Pass:   "123456",
		DBName: "v2sleep",
	}
)

func TestMain(t *testing.M) {
	Init(c)
	os.Exit(t.Run())
}
