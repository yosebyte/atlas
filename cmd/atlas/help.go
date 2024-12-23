package main

import (
	"runtime"

	"github.com/yosebyte/passport/pkg/log"
)

func helpInfo() {
	log.Info(`Version: %v %v/%v

Usage:
    atlas <core_mode>://<link_addr>#<access_addr>

Examples:
    # Run as server
    atlas server://hostname:10101

    # Run as client
    atlas client://hostname:10101#127.0.0.1:8080

Arguments:
    <core_mode>    Choose from "server" and "client" core
    <link_addr>    Interlink address to listen or connect
    <access_addr>  Optional HTTP proxy address for access

`, version, runtime.GOOS, runtime.GOARCH)
}
