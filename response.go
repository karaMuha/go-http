package http

import "strconv"

type HttpResponse struct {
	cookies    []*Cookie
	headers    map[string]string
	statusCode string
	Body       []byte
}

func (res *HttpResponse) SetCookie(cookie *Cookie) {
	res.cookies = append(res.cookies, cookie)
}

func (res *HttpResponse) SetHeader(key string, val string) {
	res.headers[key] = val + "\r\n"
}

func (res *HttpResponse) WriteStatus(code int) {
	if code < 100 || code > 599 {
		code = 200
	}
	res.statusCode = strconv.Itoa(code)
}

func (res *HttpResponse) Write(data []byte) {
	res.Body = data
}
