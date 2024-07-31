package app

import (
	"log"
	"net/http"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

func (core *Core) AccountDeleteHandler(rw http.ResponseWriter, r *http.Request) {

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

	// read user data from database
	var b []byte
	if len(value) > 0 {
		if b, err = core.Store.Read("users", value); err != nil {
			log.Printf("error read from database: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	} else {
		log.Printf("not found session, ip %v\n", r.RemoteAddr)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	// unmarshal data to model if we have user in database
	md := new(models.User)
	if len(b) > 0 {
		if err = md.UnmarshalBinary(b); err != nil {
			log.Printf("error unmarshal model: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	if md.Root {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	if err = core.Store.Delete("users", value); err != nil {
		log.Printf("error delete user from database: %v\n", err.Error())
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
