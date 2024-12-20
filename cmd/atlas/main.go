package main

import (
	"net/url"
	"os"

	"github.com/yosebyte/passport/pkg/log"
	"github.com/yosebyte/passport/pkg/tls"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		helpInfo()
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Unable to parse raw URL: %v", err)
	}
	tlsConfig, err := acme(parsedURL.Host)
	if err != nil {
		log.Error("Unable to obtain TLS config: %v", err)
		if tlsConfig, err = tls.NewTLSconfig("yosebyte/atlas:" + version); err != nil {
			log.Fatal("Unable to generate TLS config: %v", err)
		}
	}
	coreSelect(parsedURL, tlsConfig)
}
