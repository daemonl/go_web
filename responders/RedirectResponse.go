package responders

import (
	"net/http"

	"github.com/daemonl/go_web/router"
)

type RedirectResponse struct {
	URL  string
	Code int
}

func Redirect(url string) *RedirectResponse {
	return &RedirectResponse{
		URL:  url,
		Code: http.StatusSeeOther,
	}
}

func (resp *RedirectResponse) Respond(b *router.Baton) error {
	w, r := b.Raw()
	http.Redirect(w, r, resp.URL, resp.Code)
	return nil
}
