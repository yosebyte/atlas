package main

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/yosebyte/passport/pkg/conn"
	"github.com/yosebyte/passport/pkg/log"
)

func client(parsedURL *url.URL) error {
	serverAddr := parsedURL.Host
	proxyAddr := parsedURL.Fragment
	if proxyAddr == "" {
		_, port, err := net.SplitHostPort(serverAddr)
		if err != nil {
			return err
		}
		proxyAddr = net.JoinHostPort("127.0.0.1", port)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleClientRequest(w, r, serverAddr)
	})
	return http.ListenAndServe(proxyAddr, nil)
}

func handleClientRequest(w http.ResponseWriter, r *http.Request, serverAddr string) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tlsConfig := &tls.Config{}
	serverConn, err := tls.Dial("tcp", serverAddr, tlsConfig)
	if err != nil {
		log.Warn("Unable to verify server certificate: %v", err)
		tlsConfig.InsecureSkipVerify = true
		serverConn, err = tls.Dial("tcp", serverAddr, tlsConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			clientConn.Close()
			return
		}
	}
	if err := r.Write(serverConn); err != nil {
		log.Error("Unable to write request: %v", err)
		clientConn.Close()
		serverConn.Close()
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
