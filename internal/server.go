package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/yosebyte/x/conn"
	"github.com/yosebyte/x/log"
)

func NewServer(parsedURL *url.URL, tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	return &http.Server{
		Addr:      getAccessAddr(strings.TrimPrefix(parsedURL.Path, "/")),
		ErrorLog:  logger.StdLogger(),
		Handler:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serverConnect(w, r, logger) }),
		TLSConfig: tlsConfig,
	}
}

func serverConnect(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	if r.Method == http.MethodConnect {
		clientConn, err := hijackConnection(w)
		if err != nil {
			logger.Error("Hijack failed: %v", err)
			return
		}
		logger.Debug("Client connection: %v <-> %v", clientConn.LocalAddr(), clientConn.RemoteAddr())
		defer func() {
			if clientConn != nil {
				clientConn.Close()
			}
		}()
		logger.Debug("Targeting server: %v", r.URL.Host)
		targetConn, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			logger.Error("Dial failed: %v", err)
			return
		}
		logger.Debug("Target connection: %v <-> %v", targetConn.LocalAddr(), targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		if _, err := clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
			logger.Error("Write failed: %v", err)
			return
		}
		logger.Debug("Starting exchange: %v <-> %v", clientConn.LocalAddr(), targetConn.LocalAddr())
		_, _, err = conn.DataExchange(clientConn, targetConn)
		logger.Debug("Exchange complete: %v", err)
	} else {
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = r.Host
				req.RequestURI = ""
				req.Header.Del("Proxy-Connection")
			},
			ErrorLog: logger.StdLogger(),
		}
		logger.Debug("HTTP: %v -> %v", r.Method, r.URL)
		proxy.ServeHTTP(w, r)
	}
}
