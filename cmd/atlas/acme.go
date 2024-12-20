// This file uses the certmagic package, which is licensed under the Apache License 2.0.
// For more information, see https://github.com/caddyserver/certmagic/blob/master/LICENSE

package main

import (
	"crypto/tls"

	"github.com/caddyserver/certmagic"
)

func acme(hostname string) (*tls.Config, error) {
	certmagic.Default.Storage = &certmagic.FileStorage{Path: "."}
	certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "atlas@" + hostname
	tlsConfig, err := certmagic.TLS([]string{hostname})
	return tlsConfig, err
}
