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
	/* responseString := "HTTP/1.1"
	// set status code
	responseString = responseString + " " + res.statusCode + "\r\n" */
	// set cookies
	for _, cookie := range res.cookies {
		sb.WriteString(fmt.Sprintf("Set-Cookie: %s=%s", cookie.Name, cookie.Value))
		//cookieString := "Set-Cookie: " + cookie.Name + "=" + cookie.Value
		if cookie.HttpOnly {
			sb.WriteString("; HttpOnly")
			//cookieString = cookieString + "; HttpOnly"
		}
		if cookie.Secure {
			sb.WriteString("; Secure")
			//cookieString = cookieString + "; Secure"
		}
		if !cookie.Expires.IsZero() {
			sb.WriteString(fmt.Sprintf("; Expires=%s", cookie.Expires.String()))
			//cookieString = cookieString + "; Expires=" + cookie.Expires.String()
		}
		sb.WriteString("\r\n")
		/* cookieString = cookieString + "\r\n"
		responseString = responseString + cookieString */
	}
	// set headers
	for k, v := range res.headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		/* headerString := k + ": " + v + "\r\n"
		responseString = responseString + headerString */
	}

	sb.WriteString("\r\n")
	//responseString = responseString + "\r\n"

	return sb.String()
}
