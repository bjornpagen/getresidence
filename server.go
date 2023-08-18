package getresidence

import (
	"net/http"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
	"github.com/hako/branca"
	"github.com/stripe/stripe-go/v75/client"
	"github.com/uptrace/bunrouter"
)

type server struct {
	db     grDb
	branca *branca.Branca
	stripe *client.API
}

func New(branca *branca.Branca, stripe *client.API, dbUrl string) (*server, error) {
	db, err := newDb(dbUrl)
	if err != nil {
		return nil, errors.New("new db")
	}

	s := &server{
		db,
		branca,
		stripe,
	}

	return s, nil
}

func errorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		if err := next(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
}

func (s *server) Routes() *bunrouter.Router {
	r := bunrouter.New(bunrouter.Use(errorHandler))

	r.GET("/", s.getRoot)
	r.GET("/privacy", s.getPrivacy)
	r.GET("/dubai", s.getDubai())
	r.PUT("/dubai/name", s.putDubaiName())
	r.PUT("/dubai/phone", s.putDubaiPhone())
	r.PUT("/dubai/email", s.putDubaiEmail())
	r.POST("/dubai/checkout", s.postDubaiCheckout())

	return r
}

func (s *server) Close() (err error) {
	errcapture.Do(&err, s.db.Close, "close db")

	return
}
