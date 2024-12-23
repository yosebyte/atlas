package main

import (
	"runtime"

	"github.com/yosebyte/passport/pkg/log"
)

func helpInfo() {
	log.Info(`Version: %v %v/%v

Usage: atlas <core_mode>://<link_addr>/<access_addr>
`, version, runtime.GOOS, runtime.GOARCH)
}
