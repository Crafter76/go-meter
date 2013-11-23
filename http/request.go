package http

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
)

type Request struct {
	Method string

	URL *url.URL

	Header map[string][]string

	Body          io.Reader
	ContentLength int64

	Host string

	BufferSize int64
}

func (req *Request) Write(w io.Writer) error {
	headers := "Host: " + req.Host + "\r\n"
	if req.Method == "POST" || req.Method == "PUT" {
		headers += fmt.Sprintf("Content-Length: %d\r\n", req.ContentLength)
	}
	if req.Header != nil {
		for key, values := range req.Header {
			if key == "Content-Length" || key == "Host" {
				continue
			}
			for _, value := range values {
				headers += key + ": " + value + "\r\n"
			}
		}
	}
	pocket := fmt.Sprintf("%s %s HTTP/1.1\r\n%s\r\n",
		valueOrDefault(req.Method, "GET"),
		req.URL.RequestURI(),
		headers,
	)
	_, err := io.WriteString(w, pocket)
	if err != nil {
		return err
	}
	if req.Method == "POST" || req.Method == "PUT" {
		req.BufferSize = req.ContentLength
		bodyReader := bufio.NewReader(req.Body)
		_, err := bodyReader.WriteTo(w)
		if err != nil {
			return err
		}
	}
	req.BufferSize += int64(len(pocket))
	return nil
}

func valueOrDefault(value string, def string) string {
	if len(value) == 0 {
		return def
	} else {
		return value
	}
}
