package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
)

func Client(parsedURL *url.URL) error {
	serverAddr := parsedURL.Host
	accessAddr := parsedURL.Fragment
	if accessAddr == "" {
		_, port, err := net.SplitHostPort(serverAddr)
		if err != nil {
			return err
		}
		accessAddr = net.JoinHostPort("127.0.0.1", port)
	}
	server := &http.Server{
		Addr:     accessAddr,
		ErrorLog: log.NewLogger(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleClientRequest(w, r, serverAddr)
		}),
	}
	log.Info("Starting HTTP server on %v", accessAddr)
	return server.ListenAndServe()
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, serverAddr string) {
	if r.Method != http.MethodConnect {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed: %v", r.Method)
		return
	} else {
		statusOK(w)
	}
	r.Header.Set("User-Agent", getagentID())
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
	serverConn, err := tls.Dial("tcp", serverAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Error("Unable to dial TLS server: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Info("Server connected: %v", serverConn.RemoteAddr())
	defer func() {
		if serverConn != nil {
			serverConn.Close()
		}
	}()
	if err := r.Write(serverConn); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		log.Error("Unable to write request to server: %v", err)
		return
	}
	log.Info("Connection established: %v <-> %v", clientConn.RemoteAddr(), serverConn.RemoteAddr())
	if err := io.DataExchange(clientConn, serverConn); err != nil {
		log.Info("Connection closed: %v", err)
	}
}
