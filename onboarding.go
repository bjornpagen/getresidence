package getresidence

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/efficientgo/core/errors"

	"github.com/nyaruka/phonenumbers"

	"github.com/uptrace/bunrouter"
)

func (s *server) putDubaiName() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		name := r.FormValue("name")

		var e entry
		if err := validateName(name); err != nil {
			e = entry{
				Name:     "name",
				Value:    name,
				Schema:   "text",
				State:    "invalid",
				Small:    fmt.Sprintf("❌ %s", err),
				Endpoint: "/dubai/name",
			}
		} else {
			session, err := getSession(r)
			if err != nil {
				return errors.Wrap(err, "get session")
			}

			idstr, err := s.branca.DecodeToString(session)
			if err != nil {
				return errors.Wrap(err, "decode session")
			}

			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse id")
			}

			s.db.setName(id, r.FormValue("name"))

			e = entry{
				Name:     "name",
				Value:    name,
				Schema:   "text",
				State:    "valid",
				Small:    fmt.Sprintf("✅ name has been saved—%s", name),
				Endpoint: "/dubai/name",
			}
		}

		writeDefaultHeaders(w.Header(), "")
		e.Render(w)

		return nil
	}
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) > 100 {
		return errors.New("name cannot be longer than 100 characters")
	}
	return nil
}

func (s *server) putDubaiPhone() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) (err error) {
		phone := r.FormValue("phone")

		var e entry
		if validated, err := validatePhone(phone); err != nil {
			e = entry{
				Name:     "phone",
				Value:    phone,
				Schema:   "tel",
				State:    "invalid",
				Small:    fmt.Sprintf("❌ %s", err),
				Endpoint: "/dubai/phone",
			}
		} else {
			session, err := getSession(r)
			if err != nil {
				return errors.Wrap(err, "get session")
			}

			idstr, err := s.branca.DecodeToString(session)
			if err != nil {
				return errors.Wrap(err, "decode session")
			}

			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse id")
			}

			s.db.setPhone(id, validated)

			e = entry{
				Name:     "phone",
				Value:    validated,
				Schema:   "tel",
				State:    "valid",
				Small:    fmt.Sprintf("✅ phone has been saved—%s", phone),
				Endpoint: "/dubai/phone",
			}
		}

		writeDefaultHeaders(w.Header(), "")
		e.Render(w)

		return nil
	}
}

func validatePhone(phone string) (string, error) {
	if phone == "" {
		return "", errors.New("phone cannot be empty")
	}
	if len(phone) > 100 {
		return "", errors.New("phone cannot be longer than 100 characters")
	}
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL), nil
}

func (s *server) putDubaiEmail() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		email := r.FormValue("email")

		var e entry
		if err := validateEmail(email); err != nil {
			e = entry{
				Name:     "email",
				Value:    email,
				Schema:   "email",
				State:    "invalid",
				Small:    fmt.Sprintf("❌ %s", err),
				Endpoint: "/dubai/email",
			}
		} else {
			session, err := getSession(r)
			if err != nil {
				return errors.Wrap(err, "get session")
			}

			idstr, err := s.branca.DecodeToString(session)
			if err != nil {
				return errors.Wrap(err, "decode session")
			}

			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse id")
			}

			s.db.setEmail(id, r.FormValue("email"))

			e = entry{
				Name:     "email",
				Value:    email,
				Schema:   "email",
				State:    "valid",
				Small:    fmt.Sprintf("✅ email has been saved—%s", email),
				Endpoint: "/dubai/email",
			}
		}

		writeDefaultHeaders(w.Header(), "")
		e.Render(w)

		return nil
	}
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if len(email) > 100 {
		return errors.New("email cannot be longer than 100 characters")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
