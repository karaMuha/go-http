package http

import (
	"fmt"
	"strconv"
	"strings"
)

type HttpResponse struct {
	cookies    []*Cookie
	headers    map[string]string
	statusCode string
	Body       []byte
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		cookies:    make([]*Cookie, 0),
		headers:    make(map[string]string),
		Body:       make([]byte, 0),
		statusCode: "200",
	}
}

func (res *HttpResponse) SetCookie(cookie *Cookie) {
	res.cookies = append(res.cookies, cookie)
}

func (res *HttpResponse) SetHeader(key string, val string) {
	res.headers[key] = val + "\r\n"
}

func (res *HttpResponse) WriteStatus(code int) {
	if code < 100 || code > 599 {
		panic("Unexpected status code: " + strconv.Itoa(code))
	}
	res.statusCode = strconv.Itoa(code)
}

func (res *HttpResponse) WriteBody(data []byte) {
	res.Body = data
}

func (res *HttpResponse) assembleResponseString() string {
	var sb strings.Builder
	sb.WriteString("HTTP/1.1 ")
	sb.WriteString(res.statusCode)
	sb.WriteString("\r\n")
	// set cookies
	for _, cookie := range res.cookies {
		sb.WriteString(fmt.Sprintf("Set-Cookie: %s=%s", cookie.Name, cookie.Value))
		if cookie.HttpOnly {
			sb.WriteString("; HttpOnly")
		}
		if cookie.Secure {
			sb.WriteString("; Secure")
		}
		if !cookie.Expires.IsZero() {
			sb.WriteString(fmt.Sprintf("; Expires=%s", cookie.Expires.String()))
		}
		sb.WriteString("\r\n")
	}
	// set headers
	for k, v := range res.headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	sb.WriteString("\r\n")
	if len(res.Body) > 0 {
		sb.WriteString(string(res.Body))
	}

	return sb.String()
}
