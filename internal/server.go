package internal

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
	"github.com/yosebyte/x/tls"
)

func NewServer(parsedURL *url.URL) *http.Server {
	tlsConfig, err := tls.NewTLSconfig(getagentID())
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	return &http.Server{
		Addr:      parsedURL.Host,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
		TLSConfig: tlsConfig,
	}
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		if r.Header.Get("User-Agent") != getagentID() {
			statusOK(w)
		}
		clientConn, err := hijackConnection(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to hijack connection: %v", err)
			return
		}
		log.Info("Client connected: %v", clientConn.RemoteAddr())
		defer func() {
			if clientConn != nil {
				clientConn.Close()
			}
		}()
		targetConn, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			log.Error("Unable to dial target: %v", err)
			return
		}
		log.Info("Target connected: %v", targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		log.Info("Connection established: %v <-> %v", clientConn.RemoteAddr(), targetConn.RemoteAddr())
		if err := io.DataExchange(clientConn, targetConn); err != nil {
			log.Info("Connection closed: %v", err)
		}
	} else {
		if r.Header.Get("User-Agent") != getagentID() {
			statusOK(w)
			log.Warn("Invalid request: %v", r.RemoteAddr)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(r.URL)
		proxy.ServeHTTP(w, r)
	}
}
