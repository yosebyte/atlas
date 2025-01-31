package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/yosebyte/x/io"
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
		logger.Debug("Client connected: %v", r.RemoteAddr)
		clientConn, err := hijackConnection(w)
		if err != nil {
			logger.Error("Unable to hijack connection: %v", err)
			return
		}
		logger.Debug("Client hijacked: %v", clientConn.RemoteAddr())
		defer func() {
			if clientConn != nil {
				clientConn.Close()
			}
		}()
		if !strings.Contains(r.URL.Host, ":") {
			r.URL.Host += ":443"
		}
		logger.Debug("Connecting to target: %v", r.URL.Host)
		targetConn, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			logger.Error("Unable to dial target: %v", err)
			return
		}
		logger.Debug("Target connected: %v", targetConn.RemoteAddr())
		defer func() {
			if targetConn != nil {
				targetConn.Close()
			}
		}()
		_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			logger.Error("Unable to write response to client: %v", err)
			return
		}
		logger.Debug("Connection established: %v <-> %v", clientConn.RemoteAddr(), targetConn.RemoteAddr())
		if err := io.DataExchange(clientConn, targetConn); err != nil {
			logger.Debug("Connection closed: %v", err)
		}
	} else {
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   r.Host,
		})
		r.URL.Scheme = "http"
		r.URL.Host = r.Host
		r.RequestURI = ""
		r.Header.Del("Proxy-Connection")
		logger.Debug("Proxying HTTP: %v %v", r.Method, r.URL)
		proxy.ServeHTTP(w, r)
	}
}
