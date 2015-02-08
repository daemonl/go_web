package router

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

func GetRouter(auth Authenticator) (*Router, error) {
	h := &Router{
		routes:        make([]*route, 0, 0),
		authenticator: auth,
	}
	return h, nil
}

type Router struct {
	routes        []*route
	authenticator Authenticator
	Fallback      http.Handler
}

func (h *Router) AddRoute(format string, handler func(*Baton) (Responder, error)) error {
	reStr := format
	reStr = strings.Replace(reStr, "%d", "[0-9]+", -1)
	reStr = strings.Replace(reStr, "%s", "[0-9A-Za-z_]+", -1)

	re, err := regexp.Compile("^" + reStr + "$")
	if err != nil {
		return err
	}

	r := &route{
		format:  format,
		regexp:  re,
		handler: handler,
	}
	h.routes = append(h.routes, r)
	return nil
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	token, err := h.authenticator.Authenticate(w, r)
	if err != nil {
		log.Printf("Error checking cookie token: %s\n", err.Error())
		w.Write([]byte("Unknown Error"))
		w.WriteHeader(500)
		return
	}
	if token == nil {
		log.Printf("Insecure path, no token, redirect")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	baton := &Baton{
		w:     w,
		r:     r,
		path:  r.URL.Path,
		Token: token,
	}

	for _, route := range h.routes {
		if route.regexp.MatchString(baton.path) {
			baton.route = route
			resp, err := route.handler(baton)
			if err != nil {
				baton.SendError(err, http.StatusInternalServerError)
				return
			}
			err = resp.Respond(baton)
			if err != nil {
				baton.SendError(err, http.StatusInternalServerError)
			}
			return
		}
	}

	h.Fallback.ServeHTTP(w, r)
}
