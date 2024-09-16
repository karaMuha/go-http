package http

import "time"

type Cookie struct {
	Name     string
	Value    string
	Secure   bool
	HttpOnly bool
	Expires  time.Time
}
