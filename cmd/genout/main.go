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
	tagPrefixMap = map[rune]string{
		'港': "hk",
		'新': "sg",
		'台': "tw",
		'日': "jp",
		'美': "us",
		'韩': "kr",
		'加': "ca",
		'泰': "th",
		'英': "gb",
		'德': "de",
		'俄': "ru",
		'荷': "nl",
		'印': "in",
		'法': "fr",
		'阿': "ar",
		'巴': "br",
		'土': "tr",
		'澳': "au",
		'马': "my",
		'菲': "ph",
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
			if strings.ContainsRune(out.Tag, key) {
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
