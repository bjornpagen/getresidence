package getresidence

import (
	"io"
	"net/http"

	"github.com/bjornpagen/getresidence-prototype/pug/gen"
	"github.com/uptrace/bunrouter"
)

type privacy struct {
}

func (p privacy) Render(w io.Writer) {
	gen.Jade_privacy(w)
}

func (s *server) getPrivacy(w http.ResponseWriter, r bunrouter.Request) error {
	nonce := nonce()
	layout := layout{
		Main:        privacy{},
		Domain:      "getresidence.org",
		Icon:        emojiToFavicon("🔒"),
		Title:       "Privacy Policy",
		Description: "Your privacy is important to us. It is getresidence.org's policy to respect your privacy.",
		Nonce:       nonce,
	}

	writeDefaultHeaders(w.Header(), nonce)
	layout.Render(w)
	return nil
}
