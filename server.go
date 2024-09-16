package http

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	CRLF = "\r\n"
)

type Server struct {
	Port   string
	Routes Router
}

func NewServer(port string, routes Router) Server {
	return Server{
		Port:   port,
		Routes: routes,
	}
}

func (s *Server) Listen() error {
	address := "0.0.0.0:" + s.Port
	l, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go s.processRequest(conn)
	}
}

func (s *Server) processRequest(conn net.Conn) {
	defer conn.Close()

	request, err := parseRequest(conn)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	endpoint := request.Method + " " + request.URI
	if handler, ok := s.Routes.routes[endpoint]; ok {
		var response HttpResponse
		handler(&response, request)
		responseString := writeResponseString(&response)
		conn.Write([]byte(responseString))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func parseRequest(conn net.Conn) (*HttpRequest, error) {
	lines := make(chan []byte)
	go byteReader(lines, conn)
	requestLineValues := bytes.Split(<-lines, []byte(" "))

	if len(requestLineValues) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method := string(requestLineValues[0])
	target := string(requestLineValues[1])
	httpVersion := string(requestLineValues[2])

	headers := make(map[string]string)
	for {
		headerLine := <-lines
		if len(headerLine) == 0 {
			break
		}
		headerLineValues := bytes.Split(headerLine, []byte(": "))
		headers[string(headerLineValues[0])] = string(headerLineValues[1])
	}

	body := make([]byte, 0)
	for bodyLine := range lines {
		body = append(body, bodyLine...)
	}

	return &HttpRequest{
		Method:      method,
		URI:         target,
		HttpVersion: httpVersion,
		Headers:     headers,
		Body:        body,
	}, nil
}

func byteReader(channel chan []byte, connection net.Conn) {
	defer close(channel)
	buffer := make([]byte, 0)

	for {
		tmp := make([]byte, 1024)
		n, err := connection.Read(tmp)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		buffer = append(buffer, tmp[:n]...)
		splitBuffer := bytes.Split(buffer, []byte(CRLF))

		for _, line := range splitBuffer {
			channel <- line
		}

		if n <= len(tmp) || err == io.EOF {
			return
		}
	}
}

func writeResponseString(response *HttpResponse) string {
	responseString := "HTTP/1.1"
	// set status code
	responseString = responseString + " " + response.statusCode + " OK\r\n"
	// set cookies
	for _, cookie := range response.cookies {
		cookieString := "Set-Cookie: " + cookie.Name + "=" + cookie.Value
		if !cookie.Expires.IsZero() {
			cookieString = cookieString + "Expires=" + cookie.Expires.String()
		}
		if cookie.HttpOnly {
			cookieString = cookieString + "HttpOnly"
		}
		if cookie.Secure {
			cookieString = cookieString + "Secure"
		}
		cookieString = cookieString + "\r\n"
		responseString = responseString + cookieString
	}
	// set headers
	for k, v := range response.headers {
		headerString := k + ": " + v + "\r\n"
		responseString = responseString + headerString
	}

	responseString = responseString + "\r\n"

	return responseString
}
