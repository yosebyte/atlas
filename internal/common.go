package internal

import (
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func getUserAgent(userAgentName string) string {
	userAgentID := strconv.FormatInt(time.Now().Truncate(time.Minute).Unix(), 16)
	return userAgentName + "/" + userAgentID
}

func getAccessAddr(accessAddr string) string {
	if accessAddr == "" {
		return generateAccessAddr()
	}
	return accessAddr
}

func generateAccessAddr() string {
	port := rand.Intn(7169) + 1024
	return "127.0.0.1:" + strconv.Itoa(port)
}

func hijackConnection(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, _, err := hijacker.Hijack()
	return conn, err
}
