package app

import (
	"log"
	"net/http"
)

func (core *Core) AccountExitHandler(rw http.ResponseWriter, r *http.Request) {

	// check cookie
	ck, err := r.Cookie("2miner-session")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var value string
	if err = core.Secure.Decode("2miner-session", ck.Value, &value); err != nil {
		log.Printf("error decode cookie: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var cookie = http.Cookie{
		Path:     "/",
		Name:     "2miner-session",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: 3,
		MaxAge:   -1,
	}

	http.SetCookie(rw, &cookie)
	rw.WriteHeader(http.StatusOK)
	rw.Write(core.GoodRequest)
	return
}
