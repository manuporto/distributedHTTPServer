package httpserver

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/manuporto/distributedHTTPServer/pkg/httpparser"
	"github.com/manuporto/distributedHTTPServer/pkg/server"
)

type HttpServer struct {
	s server.Server
}

func readHeader(r *bufio.Reader) (string, error) {
	b, err := r.ReadString('\n')
	var header strings.Builder
	header.WriteString(b)
	for {
		if err != nil {
			return "", err
		}

		if b == "\r\n" {
			return header.String(), nil
		}
		b, err = r.ReadString('\n')
		header.WriteString(b)
	}
}

func readBody(r *bufio.Reader, size int) ([]byte, error) {
	body := make([]byte, size)
	_, err := r.Read(body)
	return body, err
}

func ReadRequest(c net.Conn) (*httpparser.HttpFrame, error) {
	r := bufio.NewReader(c)
	header, err := readHeader(r)
	fmt.Print(header)
	if err != nil {
		return nil, err
	}
	httpheader := httpparser.GetHeader(header)
	body, err := readBody(r, httpheader.ContentLength)
	return &httpparser.HttpFrame{Header: httpheader, Body: body}, err
}

func WriteResponse(c net.Conn, res *httpparser.HttpResponse) {
	fmt.Println(*res)
	var stringRes strings.Builder
	stringRes.WriteString(httpparser.Version)
	stringRes.WriteString(" ")
	stringRes.WriteString(res.Status + httpparser.Delim)

	if res.ContentType != "" {
		stringRes.WriteString(httpparser.ContentTypeResHeader + res.ContentType + httpparser.Delim)
	}
	if res.ContentLength > 0 {
		stringRes.WriteString(httpparser.ContentLengthResHeader + strconv.Itoa(res.ContentLength) + httpparser.Delim)
	}
	stringRes.WriteString(httpparser.ConnectionHeader + httpparser.Delim)
	stringRes.WriteString(httpparser.Delim)

	stringRes.WriteString(string(res.Body))

	c.Write([]byte(stringRes.String()))
}

func (hs HttpServer) ListenAndServe(address string, handler func(net.Conn)) {
	hs.s = server.Server{address, handler}
	hs.s.Serve()
}
