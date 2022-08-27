package proto

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

func NewClashSubAddrParser() *ClashSubAddrParser {
	return &ClashSubAddrParser{
		hc: &http.Client{},
	}
}

type ClashSubAddrParser struct {
	hc      *http.Client
	Proxies []*Proxy
}

func (csa *ClashSubAddrParser) Decode(data []byte) error {
	rsp, err := csa.hc.Get(string(data))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, rsp.Body)
	if err != nil {
		_ = rsp.Body.Close()
		return err
	}
	_ = rsp.Body.Close()

	decoder := yaml.NewDecoder(buf)
	cc := &ClashConfig{}
	err = decoder.Decode(cc)
	if err != nil {
		return err
	}

	csa.Proxies = append(csa.Proxies, cc.Proxies...)
	return nil
}

func (csa *ClashSubAddrParser) ToV2ray() []V2ray {
	var v2raies []V2ray
	for _, p := range csa.Proxies {
		v2raies = append(v2raies, p.ToV2rayShadowsocks())
	}
	return v2raies
}

type Proxy struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Type     string `yaml:"type"`
	Cipher   string `json:"cipher"`
	Password string `json:"password"`
}

func (p *Proxy) ToV2rayShadowsocks() *V2rayShadowsocks {
	if p.Type != "ss" {
		return nil
	}
	return &V2rayShadowsocks{
		Cipher:   p.Cipher,
		Password: p.Password,
		Server:   p.Server,
		Port:     p.Port,
		Name:     p.Name,
	}
}

type ClashConfig struct {
	Proxies []*Proxy `yaml:"proxies"`
}
