package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

func (core *Core) AccountLoginHandler(rw http.ResponseWriter, r *http.Request) {

	form := new(models.FormLoginPassword)

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var (
		v      models.User
		cookie http.Cookie
		err    error
	)

	var mp = make(map[string]models.User)

	if mp, err = core.Store.ReadAll("users"); err != nil {
		log.Printf("error read from database: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	for _, v = range mp {
		if v.Email == form.Login {
			if v.Pass == core.StringToHash(form.Password) {
				var encoded string
				if encoded, err = core.Secure.Encode("2miner-session", v.UID); err != nil {
					log.Printf("error encode cookie: %v\n", err.Error())
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write(core.BadRequest)
					return
				}
				cookie = http.Cookie{
					Path:     "/",
					Name:     "2miner-session",
					Value:    encoded,
					Secure:   true,
					HttpOnly: true,
					SameSite: 3,
					MaxAge:   31536000,
				}

				http.SetCookie(rw, &cookie)
				rw.WriteHeader(http.StatusOK)
				rw.Write(core.GoodRequest)
				return
			}
			log.Printf("error password: %v %v\n", form.Password, r.RemoteAddr)
		}
	}

	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(core.BadRequest)
	return
}
