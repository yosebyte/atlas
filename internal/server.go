package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
)

func NewServer(parsedURL *url.URL, tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	return &http.Server{
		Addr:     getAccessAddr(strings.TrimPrefix(parsedURL.Path, "/")),
		ErrorLog: logger.StdLogger(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverConnect(w, r, parsedURL.Fragment, logger)
		}),
		TLSConfig: tlsConfig,
	}
}

func serverConnect(w http.ResponseWriter, r *http.Request, userAgentName string, logger *log.Logger) {
	if r.Method == http.MethodConnect {
		userAgent := r.Header.Get("User-Agent")
		logger.Debug("User-Agent: %v", userAgent)
		if userAgent != getUserAgent(userAgentName) {
			http.Error(w, "Pending connection", http.StatusOK)
			logger.Debug("Pending connection: %v", r.RemoteAddr)
		}
		if !strings.HasPrefix(userAgent, userAgentName) {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			logger.Warn("403 Forbidden: %v %v", r.RemoteAddr, userAgent)
			return
		}
		clientConn, err := hijackConnection(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Unable to hijack connection: %v", err)
			return
		}
		logger.Debug("Client connected: %v", clientConn.RemoteAddr())
		defer func() {
			if clientConn != nil {
				clientConn.Close()
			}
		}()
		targetConn, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			logger.Error("Unable to dial target: %v", err)
			return
		}
		logger.Debug("Target connected: %v", targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		logger.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), targetConn.RemoteAddr())
		if err := io.DataExchange(clientConn, targetConn); err != nil {
			logger.Debug("Connection closed: %v", err)
		}
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		logger.Warn("404 Not Found: %v %v", r.RemoteAddr, r.Method)
		return
	}
}
