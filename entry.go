package getresidence

import (
	"io"

	"github.com/bjornpagen/getresidence/pug/gen"
)

type entry struct {
	Name     string
	Value    string
	Schema   string
	State    string
	Small    string
	Endpoint string
}

func (e entry) Render(w io.Writer) {
	gen.Jade_entry(e, w)
}
