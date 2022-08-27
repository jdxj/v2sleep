package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/router"
)

type Config struct {
	Web router.Config `yaml:"web"`
	DB  dao.Config    `yaml:"db"`
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	c := &Config{}
	return c, decoder.Decode(c)
}
