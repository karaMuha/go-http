package http

import "fmt"

type HttpRequest struct {
	Method      string
	URI         string
	HttpVersion string
	Headers     map[string]string
	Body        []byte
	cookies     map[string]*Cookie
}

func (req *HttpRequest) Cookie(key string) (*Cookie, error) {
	if cookie, ok := req.cookies[key]; ok {
		return cookie, nil
	}
	return nil, fmt.Errorf("Cookie with key %s not present", key)
}
