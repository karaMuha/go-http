package http

import "strconv"

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
	responseString := "HTTP/1.1"
	// set status code
	responseString = responseString + " " + res.statusCode + "\r\n"
	// set cookies
	for _, cookie := range res.cookies {
		cookieString := "Set-Cookie: " + cookie.Name + "=" + cookie.Value
		if cookie.HttpOnly {
			cookieString = cookieString + "; HttpOnly"
		}
		if cookie.Secure {
			cookieString = cookieString + "; Secure"
		}
		if !cookie.Expires.IsZero() {
			cookieString = cookieString + "; Expires=" + cookie.Expires.String()
		}
		cookieString = cookieString + "\r\n"
		responseString = responseString + cookieString
	}
	// set headers
	for k, v := range res.headers {
		headerString := k + ": " + v + "\r\n"
		responseString = responseString + headerString
	}

	responseString = responseString + "\r\n"

	return responseString
}
