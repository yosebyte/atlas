package internal

import (
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
	"github.com/yosebyte/x/tls"
)

func NewServer(parsedURL *url.URL) *http.Server {
	serverAddr := parsedURL.Host
	tlsConfig, err := tls.NewTLSconfig(getagentID())
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	return &http.Server{
		Addr:      serverAddr,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
		TLSConfig: tlsConfig,
	}
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") != getagentID() {
		statusOK(w)
	}
	clientConn, err := hijackConnection(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
}
