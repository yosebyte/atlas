package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
)

func NewServer(parsedURL *url.URL, tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	port := parsedURL.Port()
	if port == "" {
		port = "443"
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleServerRequest(w, r, logger)
	})
	return &http.Server{
		Addr:      net.JoinHostPort("", port),
		ErrorLog:  logger.StdLogger(),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
}

func handleServerRequest(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	password, set := parsedURL.User.Password()
	if _, p, ok := r.BasicAuth(); !ok || p != password {
		logger.Debug("Password: %v/%v", p, ok)
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logger.Warn("Unauthorized access: %v", r.RemoteAddr)
		return
	}
	if r.Method == http.MethodConnect {
		logger.Debug("User-Agent: %v", r.Header.Get("User-Agent"))
		if r.Header.Get("User-Agent") != getagentID() {
			http.Error(w, "Connection Established", http.StatusOK)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Debug("Method not allowed: %v/%v", r.RemoteAddr, r.Method)
		return
	}
}
