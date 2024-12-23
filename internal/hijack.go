package internal

import (
	"errors"
	"net"
	"net/http"
)

func hijackConnection(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("hijacking not supported")
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
