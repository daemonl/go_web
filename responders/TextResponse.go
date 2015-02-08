package responders

import (
	"fmt"

	"github.com/daemonl/go_web/router"
)

type TextResponse struct {
	String string
}

func (resp *TextResponse) Respond(b *router.Baton) error {
	w, _ := b.Raw()
	fmt.Fprint(w, resp.String)
	return nil
}
