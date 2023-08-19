package getresidence

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bjornpagen/getresidence/pug/gen"
	"github.com/efficientgo/core/errors"

	"github.com/uptrace/bunrouter"
)

type dubai struct {
	Onboarding gen.Component
}

func (d dubai) Render(w io.Writer) {
	gen.Jade_dubai(d, w)
}

type onboarding struct {
	Children []gen.Component
}

func (o onboarding) Render(w io.Writer) {
	gen.Jade_dubaionboarding(o, w)
}

func (s *server) getDubai() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		var id int64
		session, err := getSession(r)
		if err == nil { // session cookie exists
			idstr, err := s.branca.DecodeToString(session)
			if err != nil {
				return errors.Wrap(err, "decode session")
			}

			id, err = strconv.ParseInt(idstr, 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse id")
			}
		} else { // session cookie does not exist
			id, err = s.db.newRow()
			if err != nil {
				return errors.Wrap(err, "new row")
			}

			session, err = s.branca.EncodeToString(fmt.Sprintf("%d", id))
			if err != nil {
				return errors.Wrap(err, "encode session")
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    session,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})
		}

		name, email, phone, err := s.db.getOnboarding(id)
		if err != nil {
			return errors.Wrap(err, "get onboarding")
		}

		ob := onboarding{
			Children: []gen.Component{
				entry{
					Name:     "name",
					Value:    name,
					Schema:   "text",
					State:    "",
					Small:    "",
					Endpoint: "/dubai/name",
				},
				entry{
					Name:     "email",
					Value:    email,
					Schema:   "email",
					State:    "",
					Small:    "",
					Endpoint: "/dubai/email",
				},
				entry{
					Name:     "phone",
					Value:    max(phone, "+1"),
					Schema:   "tel",
					State:    "",
					Small:    "",
					Endpoint: "/dubai/phone",
				},
			},
		}

		main := dubai{
			Onboarding: ob,
		}

		nonce := nonce()
		layout := layout{
			Main:        main,
			Domain:      "getresidence.org",
			Icon:        emojiToFavicon("🇦🇪"),
			Title:       "Get Legal Residency in Dubai",
			Description: "Get Legal Residency in Dubai. Legally pay zero Taxes, or close to it!",
			Nonce:       nonce,
		}

		writeDefaultHeaders(w.Header(), nonce)
		layout.Render(w)

		return nil
	}
}

func max(a, b string) string {
	if len(a) > len(b) {
		return a
	}
	return b
}
