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
	subAddr   = flag.String("sub-addr", "", "v2rayN(G) sub addr")
	include   = flag.String("include", "", "include keyword in tag")
	exclude   = flag.String("exclude", "", "exclude keyword in tag")
	output    = flag.String("output", "outbounds.json", "output path")
	tagPrefix = flag.String("tag-prefix", "proxy", "tag prefix")
	proto     = flag.String("proto", "vmess", "filter by proto")
)

func main() {
	flag.Parse()
	if *subAddr == "" {
		panic("empty sub addr")
	}

	sap := v2rayng.NewSubAddrParser()
	err := sap.Decode([]byte(*subAddr))
	if err != nil {
		logrus.Fatalf("decode err: %s", err)
	}

	d, err := sap.Outbounds(
		excludeKeyword(*exclude),
		includeKeyword(*include),
		filterProto(*proto),
		addTagPrefix(*tagPrefix),
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

func addTagPrefix(tagPrefix string) v2rayng.Filter {
	return func(out *v2raycore.Outbound) bool {
		if tagPrefix != "" {
			out.Tag = fmt.Sprintf("%s_%s", tagPrefix, out.Tag)
		}
		return true
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
		if p == "" {
			return true
		}
		if out.Protocol == p {
			return true
		}
		return false
	}
}
