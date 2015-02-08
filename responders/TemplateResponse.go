package responders

import (
	"fmt"

	"github.com/daemonl/go_sweetpl"
	"github.com/daemonl/go_web/router"
)

var DefaultSweeTpl *sweetpl.SweeTpl

type TemplateResponse struct {
	Data     interface{}
	Template string
	Tpl      *sweetpl.SweeTpl
}

func (resp *TemplateResponse) Respond(b *router.Baton) error {
	w, _ := b.Raw()
	if resp.Tpl == nil && DefaultSweeTpl == nil {
		return fmt.Errorf("No template engine found")
	}
	t := resp.Tpl
	if t == nil {
		t = DefaultSweeTpl
	}
	return t.Render(w, resp.Template, resp.Data)
}
