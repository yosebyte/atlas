package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/yosebyte/x/conn"
	"github.com/yosebyte/x/log"
)

func NewClient(parsedURL *url.URL, logger *log.Logger) *http.Server {
	return &http.Server{
		Addr:     getAccessAddr(strings.TrimPrefix(parsedURL.Path, "/")),
		ErrorLog: logger.StdLogger(),
		Handler:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { clientConnect(w, r, parsedURL, logger) }),
	}
}

func clientConnect(w http.ResponseWriter, r *http.Request, parsedURL *url.URL, logger *log.Logger) {
	clientConn, err := hijackConnection(w)
	if err != nil {
		logger.Error("Hijack failed: %v", err)
		return
	}
	logger.Debug("Client connection: %v <-> %v", clientConn.LocalAddr(), clientConn.RemoteAddr())
	defer func() {
		if clientConn != nil {
			clientConn.Close()
		}
	}()
	tlsConfig := &tls.Config{}
	if net.ParseIP(parsedURL.Hostname()) != nil {
		tlsConfig.InsecureSkipVerify = true
	}
	serverConn, err := tls.Dial("tcp", parsedURL.Host, tlsConfig)
	if err != nil {
		logger.Error("Dial failed: %v", err)
		return
	}
	logger.Debug("Server connection: %v <-> %v", serverConn.LocalAddr(), serverConn.RemoteAddr())
	defer func() {
		if serverConn != nil {
			serverConn.Close()
		}
	}()
	if err := r.Write(serverConn); err != nil {
		logger.Error("Write failed: %v", err)
		return
	}
	logger.Debug("Starting exchange: %v <-> %v", clientConn.LocalAddr(), serverConn.LocalAddr())
	_, _, err = conn.DataExchange(clientConn, serverConn)
	logger.Debug("Exchange complete: %v", err)
}
