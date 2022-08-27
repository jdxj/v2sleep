package dao

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type Config struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	DBName string `yaml:"db_name"`
}

func Init(conf Config) {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.DBName)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		logrus.Fatalf("open mysql err: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("get sql db err: %s", err)
	}
	if err = sqlDB.Ping(); err != nil {
		logrus.Fatalf("ping db err: %s", err)
	}

	DB = db
}
