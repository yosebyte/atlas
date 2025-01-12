package internal

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

var agentPrefix = "nghttp2/"

func hijackConnection(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, http.ErrNotSupported
	}
	conn, _, err := hijacker.Hijack()
	return conn, err
}

func getagentID() string {
	return agentPrefix + strconv.FormatInt(time.Now().Truncate(time.Minute).Unix(), 16)
}

const pageNotFound = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 Not Found</title>
    <style>
        body {
			font-family: Courier, monospace;
            text-align: center;
        }
    </style>
</head>
<body>
    <h1>404 Not Found</h1>
    <p>The page you are looking for does not exist.</p>
	<a href="/">Try again</a>
</body>
</html>
`
