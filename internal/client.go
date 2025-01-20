package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
)

func NewClient(parsedURL *url.URL, logger *log.Logger) *http.Server {
	port := parsedURL.Port()
	if port == "" {
		port = "443"
	}
	parsedURL.Host = net.JoinHostPort(parsedURL.Hostname(), port)
	accessAddr := strings.TrimPrefix(parsedURL.Path, "/")
	if accessAddr == "" {
		accessAddr = getAccessAddr()
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleClientRequest(w, r, parsedURL, logger)
	})
	return &http.Server{
		Addr:     accessAddr,
		ErrorLog: logger.StdLogger(),
		Handler:  handler,
	}
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, parsedURL *url.URL, logger *log.Logger) {
	if r.Method == http.MethodConnect {
		http.Error(w, "Pending connection", http.StatusOK)
		logger.Debug("Pending connection: %v", r.RemoteAddr)
		r.Header.Set("User-Agent", getUserAgent(parsedURL.Fragment))
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
		tlsConfig := &tls.Config{}
		if net.ParseIP(parsedURL.Hostname()) != nil {
			tlsConfig.InsecureSkipVerify = true
			logger.Debug("Skipping cert verification: %v", parsedURL.Hostname())
		}
		serverConn, err := tls.Dial("tcp", parsedURL.Host, tlsConfig)
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
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   r.Host,
			Path:   r.URL.Path,
		})
		proxy.ErrorLog = logger.StdLogger()
		logger.Debug("Handling HTTP request: %v %v", r.Method, r.URL)
		proxy.ServeHTTP(w, r)
	}
}
