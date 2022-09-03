package clash

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

func NewSubAddrParser() *SubAddrParser {
	return &SubAddrParser{
		hc: &http.Client{},
	}
}

type SubAddrParser struct {
	hc      *http.Client
	Proxies []*Proxy
}

func (csa *SubAddrParser) Decode(data []byte) error {
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
	cc := &Config{}
	err = decoder.Decode(cc)
	if err != nil {
		return err
	}

	csa.Proxies = append(csa.Proxies, cc.Proxies...)
	return nil
}

type Proxy struct {
	Name           string `yaml:"name"`
	Server         string `yaml:"server"`
	Port           int64  `yaml:"port"`
	Type           string `yaml:"type"`
	Cipher         string `json:"cipher"`
	Password       string `json:"password"`
	SNI            string `json:"sni"`
	SkipCertVerify bool   `json:"skip-cert-verify"`
	UDP            bool   `json:"udp"`
}

type Config struct {
	Proxies []*Proxy `yaml:"proxies"`
}
