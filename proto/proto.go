package proto

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"gopkg.in/yaml.v3"
)

func NewClashConfig(r io.Reader) (*ClashConfig, error) {
	decoder := yaml.NewDecoder(r)
	cc := &ClashConfig{}
	return cc, decoder.Decode(cc)
}

type Proxy struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Type     string `yaml:"type"`
	Cipher   string `json:"cipher"`
	Password string `json:"password"`
}

type ClashConfig struct {
	Proxies []Proxy `yaml:"proxies"`
}

func (cc *ClashConfig) ToV2rayShadowsocks() []*V2rayShadowsocks {
	var vsss []*V2rayShadowsocks
	for _, p := range cc.Proxies {
		vsss = append(vsss, &V2rayShadowsocks{
			Cipher:   p.Cipher,
			Password: p.Password,
			Server:   p.Server,
			Port:     p.Port,
			Name:     p.Name,
		})
	}
	return vsss
}

type V2rayShadowsocks struct {
	Cipher   string
	Password string
	Server   string
	Port     int
	Name     string
}

func (vss *V2rayShadowsocks) String() string {
	path := fmt.Sprintf("%s:%s@%s:%d", vss.Cipher, vss.Password, vss.Server, vss.Port)
	path = base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(path))
	anchor := url.PathEscape(vss.Name)
	return fmt.Sprintf("ss://%s#%s", path, anchor)
}
