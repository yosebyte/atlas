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

func NewClient(parsedURL *url.URL, logger *log.Logger) *http.Server {
	serverAddr := parsedURL.Host
	accessAddr := strings.TrimPrefix(parsedURL.Path, "/")
	if accessAddr == "" {
		_, port, err := net.SplitHostPort(serverAddr)
		if err != nil {
			logger.Error("Unable to split host and port: %v", err)
			return nil
		}
		accessAddr = net.JoinHostPort("127.0.0.1", port)
		if accessAddr == serverAddr {
			accessAddr = net.JoinHostPort("127.0.0.2", port)
		}
	}
	return &http.Server{
		Addr:     accessAddr,
		ErrorLog: logger.StdLogger(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleClientRequest(w, r, serverAddr, logger)
		}),
	}
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, serverAddr string, logger *log.Logger) {
	if r.Method == http.MethodConnect {
		http.Error(w, "Connection Established", http.StatusOK)
		r.Header.Set("User-Agent", getagentID())
		logger.Debug("User-Agent: %v", r.Header.Get("User-Agent"))
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
		serverConn, err := tls.Dial("tcp", serverAddr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			logger.Error("Unable to dial server: %v", err)
			return
		}
		logger.Debug("Server connected: %v", serverConn.RemoteAddr())
		defer func() {
			if serverConn != nil {
				serverConn.Close()
			}
		}()
		if err := r.Write(serverConn); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			logger.Error("Unable to write request to server: %v", err)
			return
		}
		logger.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), serverConn.RemoteAddr())
		if err := io.DataExchange(clientConn, serverConn); err != nil {
			logger.Debug("Connection closed: %v", err)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Debug("Method not allowed: %v/%v", r.RemoteAddr, r.Method)
		return
	}
}
