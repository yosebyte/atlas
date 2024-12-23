package internal

import (
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/passport/pkg/conn"
	"github.com/yosebyte/passport/pkg/log"
	"github.com/yosebyte/passport/pkg/tls"
)

func Server(parsedURL *url.URL) error {
	listenAddr := parsedURL.Host
	tlsConfig, err := tls.NewTLSconfig(listenAddr)
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	server := &http.Server{
		Addr:      listenAddr,
		TLSConfig: tlsConfig,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
	}
	log.Info("Starting HTTPS server on %v", listenAddr)
	return server.ListenAndServeTLS("", "")
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	clientConn, err := hijackConnection(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	defer func() {
		if targetConn != nil {
			targetConn.Close()
		}
	}()
	if _, err := w.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		log.Error("Failed to write connection established response: %v", err)
		return
	}
	if err := conn.DataExchange(clientConn, targetConn); err != nil {
		if err == io.EOF {
			log.Info("Connection closed successfully: %v", err)
		} else {
			log.Warn("Connection closed unexpectedly: %v", err)
		}
	}
}
