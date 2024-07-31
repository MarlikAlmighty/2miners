package app

import (
	"net/http"
)

func (core *Core) NotFound(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
