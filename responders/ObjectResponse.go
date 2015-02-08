package responders

import (
	"encoding/json"

	"github.com/daemonl/go_web/router"
)

type ObjectResponse struct {
	Obj  interface{}
	Code int
}

func (resp *ObjectResponse) Respond(b *router.Baton) error {
	w, _ := b.Raw()
	if resp.Code != 0 {
		w.WriteHeader(resp.Code)
	}
	enc := json.NewEncoder(w)
	return enc.Encode(resp.Obj)
}
