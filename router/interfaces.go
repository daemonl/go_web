package router

import "net/http"

type Responder interface {
	Respond(*Baton) error
}

type Authenticator interface {
	Authenticate(http.ResponseWriter, *http.Request) (interface{}, error)
}
