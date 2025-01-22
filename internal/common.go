package internal

import (
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func getUserAgent() string {
	return "atlas/" + strconv.FormatInt(time.Now().Truncate(time.Minute).Unix(), 16)
}

func getAccessAddr(accessAddr string) string {
	if accessAddr != "" {
		return accessAddr
	}
	return "127.0.0.1:" + strconv.Itoa(rand.Intn(7169)+1024)
}

func hijackConnection(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, _, err := hijacker.Hijack()
	return conn, err
}
