package main

import (
	"runtime"

	"github.com/yosebyte/x/log"
)

func helpInfo() {
	log.Info(`Version: %v %v/%v

Usage: atlas <core_mode>://<server_addr>#<access_addr>
`, version, runtime.GOOS, runtime.GOARCH)
}
