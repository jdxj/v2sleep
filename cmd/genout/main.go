package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/proto/v2raycore"
	"github.com/jdxj/v2sleep/proto/v2rayng"
)

var (
	subAddr = flag.String("sub-addr", "", "v2rayN(G) sub addr")
	include = flag.String("include", "", "include keyword in tag")
	exclude = flag.String("exclude", "", "exclude keyword in tag")
	output  = flag.String("output", "06_outbounds_sub.json", "output path")
	proto   = flag.String("proto", "", "filter by proto")
)

var (
	tagPrefixMap = map[string]string{
		"香港": "hk",
		"广港": "hk",

		"新加坡": "sg",
		"广新":   "sg",

		"台湾": "tw",
		"广台": "tw",

		"日本": "jp",
		"广日": "jp",

		"美国": "us",
		"广美": "us",

		"韩国":     "kr",
		"加拿大":   "ca",
		"泰国":     "th",
		"英国":     "gb",
		"德国":     "de",
		"俄罗斯":   "ru",
		"荷兰":     "nl",
		"印度":     "in",
		"法国":     "fr",
		"阿根廷":   "ar",
		"巴西":     "br",
		"土耳其":   "tr",
		"澳大利亚": "au",
		"马来西亚": "my",
		"菲律宾":   "ph",
	}
)

func main() {
	flag.Parse()
	if *subAddr == "" {
		panic("empty sub addr")
	}

	var (
		addrs = strings.Split(*subAddr, ",")
		sap   = v2rayng.NewSubAddrParser()
	)
	for _, addr := range addrs {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		err := sap.Decode([]byte(addr))
		if err != nil {
			logrus.Fatalf("decode err: %s", err)
		}
	}

	d, err := sap.Outbounds(
		excludeKeyword(*exclude),
		includeKeyword(*include),
		filterProto(*proto),
		distinct(),
		addTagPrefix(),
	)
	if err != nil {
		logrus.Fatalf("gen outbounds config err: %s", err)
	}

	f, err := os.Create(*output)
	if err != nil {
		logrus.Fatalf("create %s err: %s", *output, err)
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	_, err = f.Write(d)
	if err != nil {
		logrus.Fatalf("write %s err: %s", *output, err)
	}
}

func addTagPrefix() v2rayng.Filter {
	return func(out *v2raycore.Outbound) bool {
		for key, pre := range tagPrefixMap {
			if strings.Contains(out.Tag, key) {
				out.Tag = fmt.Sprintf("%s_%s", pre, out.Tag)
				return true
			}
		}
		logrus.Warnf("tag prefix not matched: %s\n", out.Tag)
		return false
	}
}

func includeKeyword(kw string) v2rayng.Filter {
	return func(out *v2raycore.Outbound) bool {
		if kw == "" {
			return true
		}
		if strings.Contains(out.Tag, kw) {
			return true
		}
		return false
	}
}

func excludeKeyword(kw string) v2rayng.Filter {
	return func(out *v2raycore.Outbound) bool {
		if kw == "" {
			return true
		}
		if strings.Contains(out.Tag, kw) {
			return false
		}
		return true
	}
}

func filterProto(p string) v2rayng.Filter {
	return func(out *v2raycore.Outbound) bool {
		if p == "" || out.Protocol == p {
			return true
		}
		return false
	}
}

func distinct() v2rayng.Filter {
	set := make(map[string]int)
	return func(out *v2raycore.Outbound) bool {
		set[out.Tag]++
		if set[out.Tag] > 1 {
			logrus.Warningf("repeated %s, count: %d, proto: %s", out.Tag, set[out.Tag], out.Protocol)
			return false
		}
		return true
	}
}
