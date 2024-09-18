package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	CRLF = "\r\n"
)

type Server struct {
	Port   string
	Routes *Router
}

func NewServer(port string, routes *Router) Server {
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
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	request, err := parseRequest(conn)
	if err != nil {
		_, err = conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	endpoint := request.Method + " " + request.URI
	if handler, ok := s.Routes.routes[endpoint]; ok {
		response := NewHttpResponse()
		handler(response, request)
		responseString := response.assembleResponseString()
		_, err = conn.Write([]byte(responseString))
		if err != nil {
			log.Println(err)
		}
		return
	}

	_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		log.Println(err)
	}
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
	cookies := make(map[string]*Cookie)
	for {
		headerLine := <-lines
		if len(headerLine) == 0 {
			break
		}
		headerLineValues := bytes.Split(headerLine, []byte(": "))
		if string(headerLineValues[0]) == "Cookie" {
			cookies = parseCookies(headerLineValues[1])
			continue
		}
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
		cookies:     cookies,
	}, nil
}

func byteReader(channel chan []byte, connection net.Conn) {
	defer close(channel)
	buffer := make([]byte, 0)

	for {
		tmp := make([]byte, 1024)
		n, err := connection.Read(tmp)
		if err != nil && err != io.EOF {
			log.Println(err)
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

func parseCookies(cookiebytes []byte) map[string]*Cookie {
	cookies := make(map[string]*Cookie)
	cookieArray := bytes.Split(cookiebytes, []byte("; "))
	for _, v := range cookieArray {
		cookieSplit := bytes.Split(v, []byte("="))
		cookieName := cookieSplit[0]
		cookieValue := cookieSplit[1]
		cookie := &Cookie{
			Name:  string(cookieName),
			Value: string(cookieValue),
		}
		cookies[string(cookieName)] = cookie
	}

	return cookies
}
