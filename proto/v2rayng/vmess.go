package v2rayng

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jdxj/v2sleep/proto/v2raycore"
)

type Vmess struct {
	TagPrefix string `json:"-"`

	Version   json.RawMessage `json:"v"`
	Name      string          `json:"ps"`
	Address   string          `json:"add"`
	Port      json.RawMessage `json:"port"`
	ID        string          `json:"id"`
	AID       json.RawMessage `json:"aid"`
	Security  string          `json:"scy"`
	TransType string          `json:"net"`
	FakeType  string          `json:"type"`
	Host      string          `json:"host"`
	Path      string          `json:"path"`
	TLS       string          `json:"tls"`
	SNI       string          `json:"sni"`
}

func (vv *Vmess) Encode() ([]byte, error) {
	data, err := json.Marshal(vv)
	if err != nil {
		return nil, err
	}
	s := base64.StdEncoding.EncodeToString(data)
	return []byte(fmt.Sprintf("vmess://%s", s)), nil
}

func (vv *Vmess) Decode(data []byte) error {
	s := strings.TrimPrefix(string(data), "vmess://")
	jsonData, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, vv)
}

func (vv *Vmess) Outbound() (*v2raycore.Outbound, error) {
	out := &v2raycore.Outbound{
		Tag:      vv.Name,
		Protocol: "vmess",
		Settings: v2raycore.VmessSettings{
			VNext: []v2raycore.VNext{
				{
					Address: vv.Address,
					Port:    JsonRawToInt(vv.Port),
					Users: []v2raycore.User{
						{
							ID:      vv.ID,
							AlterId: JsonRawToInt(vv.AID),
							Security: func() string {
								if vv.Security != "" {
									return vv.Security
								}
								return "auto"
							}(),
						},
					},
				},
			},
		},
	}

	ss := &v2raycore.StreamSettings{
		Network: vv.TransType,
	}

	if vv.TLS == "tls" {
		ss.Security = vv.TLS
		ss.TlsSettings = &v2raycore.TlsSettings{
			ServerName:    vv.SNI,
			AllowInsecure: false,
		}
	} else {
		ss.Security = "none"
	}

	switch vv.TransType {
	case "tcp":
		ss.TcpSettings = &v2raycore.TcpSettings{
			Header: v2raycore.HttpHeader{Type: vv.FakeType},
		}
		switch vv.FakeType {
		case "none":
		case "http":
			ss.TcpSettings.Header.Request = &v2raycore.Request{
				Version: "1.1",
				Method:  "GET",
				Path:    []string{vv.Path},
				Headers: map[string][]string{
					"Host": {vv.Host},
				},
			}
		default:
			return nil, fmt.Errorf("fake type not implement: %s", vv.FakeType)
		}
	case "http", "h2":
		ss.HttpSettings = &v2raycore.HttpSettings{
			Host:   []string{vv.Host},
			Path:   vv.Path,
			Method: "GET",
		}
	case "ws":
		ss.WSSettings = &v2raycore.WSSettings{
			Path: vv.Path,
		}
	default:
		return nil, fmt.Errorf("trans type %s not implement: %s", vv.TransType, vv.Name)
	}

	out.StreamSettings = ss
	return out, nil
}

func JsonRawToInt(msg json.RawMessage) int64 {
	var i any
	err := json.Unmarshal(msg, &i)
	if err != nil {
		panic(err)
	}

	switch port := i.(type) {
	case float64:
		return int64(port)
	case string:
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			panic(err)
		}
		return p
	default:
		panic(fmt.Errorf("can not parse port: %+v", msg))
	}
}
