package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/x/io"
	"github.com/yosebyte/x/log"
)

func RunClient(parsedURL *url.URL) error {
	serverAddr := parsedURL.Host
	accessAddr := parsedURL.Fragment
	if accessAddr == "" {
		_, port, err := net.SplitHostPort(serverAddr)
		if err != nil {
			log.Error("Unable to split host and port: %v", err)
			return nil
		}
		accessAddr = net.JoinHostPort("127.0.0.1", port)
	}
	client := &http.Server{
		Addr:     accessAddr,
		ErrorLog: log.NewLogger(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleClientRequest(w, r, serverAddr)
		}),
	}
	return client.ListenAndServe()
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, serverAddr string) {
	if r.Method == http.MethodConnect {
		statusOK(w)
	}
	r.Header.Set("User-Agent", getagentID())
	log.Debug("User-Agent: %v", r.Header.Get("User-Agent"))
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
	serverConn, err := tls.Dial("tcp", serverAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		log.Error("Unable to dial server: %v", err)
		return
	}
	log.Debug("Server connected: %v", serverConn.RemoteAddr())
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
	log.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), serverConn.RemoteAddr())
	if err := io.DataExchange(clientConn, serverConn); err != nil {
		log.Debug("Connection closed: %v", err)
	}
}
