package model

import (
	"context"
	"os"
	"testing"

	"github.com/jdxj/v2sleep/dao"
)

var (
	daoC = dao.Config{
		Host:   "192.168.50.200",
		Port:   3306,
		User:   "root",
		Pass:   "123456",
		DBName: "v2sleep",
	}
)

func TestMain(t *testing.M) {
	dao.Init(daoC)
	os.Exit(t.Run())
}

func TestDeleteSubConfig(t *testing.T) {
	err := DeleteSubConfig(context.Background(), nil)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
