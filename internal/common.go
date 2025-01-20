package internal

import (
	"fmt"
	"math/rand"
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

func getUserAgent(s string) string {
	userAgent := s
	if userAgent == "" {
		userAgent = "curl"
	}
	return userAgent + "/" + strconv.FormatInt(time.Now().Truncate(time.Minute).Unix(), 16)
}

func getAccessAddr() string {
	ip := fmt.Sprintf("127.0.0.%d", rand.Intn(255))
	port := rand.Intn(7169) + 1024
	return fmt.Sprintf("%s:%d", ip, port)
}
