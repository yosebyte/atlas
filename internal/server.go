package internal

import (
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
	"github.com/yosebyte/x/tls"
)

func RunServer(parsedURL *url.URL) error {
	serverAddr := parsedURL.Host
	tlsConfig, err := tls.NewTLSconfig(getagentID())
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	server := &http.Server{
		Addr:      serverAddr,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
		TLSConfig: tlsConfig,
	}
	return server.ListenAndServeTLS("", "")
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		log.Debug("User-Agent: %v", r.Header.Get("User-Agent"))
		if r.Header.Get("User-Agent") != getagentID() {
			statusOK(w)
		}
		clientConn, err := hijackConnection(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to hijack connection: %v", err)
			return
		}
		log.Debug("Client connected: %v", clientConn.RemoteAddr())
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
		log.Debug("Target connected: %v", targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		log.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), targetConn.RemoteAddr())
		if err := io.DataExchange(clientConn, targetConn); err != nil {
			log.Debug("Connection closed: %v", err)
		}
	} else {
		statusOK(w)
	}
}
