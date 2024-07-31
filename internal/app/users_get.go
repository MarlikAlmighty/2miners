package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

func (core *Core) UsersGetHandler(rw http.ResponseWriter, r *http.Request) {

	// Check cookie
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

	var mu = make(map[string]models.User)

	// read user data from database
	if len(value) > 0 {
		if mu, err = core.Store.ReadAll("users"); err != nil {
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

	var root bool
	var md = new(models.Users)
	for _, v := range mu {
		if v.UID == value {
			root = true
		}
		*md = append(*md, &v)
	}

	if root {
		b, _ := json.Marshal(md)
		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	}

	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(core.BadRequest)
	return
}
