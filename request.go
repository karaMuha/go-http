package http

type HttpRequest struct {
	Method      string
	URI         string
	HttpVersion string
	Headers     map[string]string
	Body        []byte
}
