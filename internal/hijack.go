package internal

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

func hijackConnection(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, _, err := hijacker.Hijack()
	return conn, err
}

func gethijackID() string {
	return strconv.FormatInt(time.Now().Truncate(time.Minute).Unix(), 16)
}
