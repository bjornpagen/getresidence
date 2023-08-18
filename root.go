package getresidence

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *server) getRoot(w http.ResponseWriter, r bunrouter.Request) error {
	// redirect to /dubai
	http.Redirect(w, r.Request, "/dubai", http.StatusFound)
	return nil
}
