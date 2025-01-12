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
	port := parsedURL.Port()
	if port == "" {
		port = "443"
	}
	serverAddr := net.JoinHostPort("", port)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleServerRequest(w, r, logger)
	})
	return &http.Server{
		Addr:      serverAddr,
		ErrorLog:  logger.StdLogger(),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
}

func handleServerRequest(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	if r.Method == http.MethodConnect {
		userAgent := r.Header.Get("User-Agent")
		logger.Debug("User-Agent: %v", userAgent)
		if userAgent != getagentID() {
			http.Error(w, "Pending connection", http.StatusOK)
			logger.Debug("Pending connection: %v", r.RemoteAddr)
		}
		if !strings.HasPrefix(userAgent, agentPrefix) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			logger.Warn("403: %v %v", r.RemoteAddr, userAgent)
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
		http.Error(w, pageNotFound, http.StatusNotFound)
		logger.Warn("404: %v %v", r.RemoteAddr, r.Method)
		return
	}
}
