package convert

import (
	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/proto"
	"github.com/jdxj/v2sleep/proto/clash"
	"github.com/jdxj/v2sleep/proto/v2rayng"
)

func ProxyToV2rayNG(proxies ...*clash.Proxy) []proto.V2rayNG {
	var v2raies []proto.V2rayNG
	for _, p := range proxies {
		var v proto.V2rayNG
		switch p.Type {
		case "ss":
			v = &v2rayng.V2rayShadowsocks{
				Cipher:   p.Cipher,
				Password: p.Password,
				Server:   p.Server,
				Port:     p.Port,
				Name:     p.Name,
			}

		case "trojan":
			v = &v2rayng.V2rayTrojan{
				Password:   p.Password,
				Server:     p.Server,
				Port:       p.Port,
				Security:   "tls",
				HeaderType: "none",
				Type:       "tcp",
				Name:       p.Name,
			}

		default:
			logrus.Warnf("%s can not to V2rayNG", p.Type)
			continue
		}
		v2raies = append(v2raies, v)
	}

	return v2raies
}
