package internal

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/passport/pkg/conn"
	"github.com/yosebyte/passport/pkg/log"
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
		Handler:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleClientRequest(w, r, serverAddr)
		}),
	}
	log.Info("Starting HTTP server on %v", accessAddr)
	return server.ListenAndServe()
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, serverAddr string) {
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
	serverConn, err := tls.Dial("tcp", serverAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Error("Unable to dial TLS server: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() {
		if serverConn != nil {
			serverConn.Close()
		}
	}()
	if _, err := w.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		log.Error("Failed to write connection established response: %v", err)
		return
	}

	if err := conn.DataExchange(clientConn, serverConn); err != nil {
		if err == io.EOF {
			log.Info("Connection closed successfully: %v", err)
		} else {
			log.Warn("Connection closed unexpectedly: %v", err)
		}
	}
}