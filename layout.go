package getresidence

import (
	_ "embed"
	"io"

	"github.com/bjornpagen/getresidence/pug/gen"
)

//go:embed static/reset.css
var RESET string

//go:embed static/style.css
var STYLE string

type layout struct {
	Main        gen.Component
	Domain      string
	Icon        string
	Title       string
	Description string
	Nonce       string
}

func (l layout) Render(w io.Writer) {
	gen.Jade_layout(l, RESET, STYLE, w)
}
