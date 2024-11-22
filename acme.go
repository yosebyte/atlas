// Copyright 2015 Matthew Holt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"

	"github.com/caddyserver/certmagic"
)

func acme(hostname string) (*tls.Config, error) {
	certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "atlas@" + hostname
	tlsConfig, err := certmagic.TLS([]string{hostname})
	return tlsConfig, err
}
