package internal

import (
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
	"github.com/yosebyte/x/tls"
)

func Server(parsedURL *url.URL) error {
	serverAddr := parsedURL.Host
	tlsConfig, err := tls.NewTLSconfig(gethijackID())
	if err != nil {
		log.Fatal("Unable to generate TLS config: %v", err)
	}
	server := &http.Server{
		Addr:      serverAddr,
		ErrorLog:  log.NewLogger(),
		Handler:   http.HandlerFunc(handleServerRequest),
		TLSConfig: tlsConfig,
	}
	log.Info("Starting HTTPS server on %v", serverAddr)
	return server.ListenAndServeTLS("", "")
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.RemoteAddr))
		w.(http.Flusher).Flush()
		log.Warn("Method not allowed: %v", r.RemoteAddr)
		return
	}
	if r.Header.Get("User-Agent") != gethijackID() {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
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
