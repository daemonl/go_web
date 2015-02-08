package router

import "regexp"

type route struct {
	format  string
	regexp  *regexp.Regexp
	handler func(b *Baton) (Responder, error)
}
