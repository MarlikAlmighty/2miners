package app

import (
	"encoding/json"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"log"
	"net/http"
)

func (core *Core) AddrDeleteHandler(rw http.ResponseWriter, r *http.Request) {

	// check cookie
	ck, err := r.Cookie("2miner-session")
	if err != nil {
		log.Printf("not found session, ip %v\n", r.RemoteAddr)
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
	if b, err = core.Store.Read("users", value); err != nil {
		log.Printf("error read from database: %v\n", err.Error())
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
	} else {
		log.Printf("not found session, ip %v\n", r.RemoteAddr)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	addr := new(models.Addr)

	if err = json.NewDecoder(r.Body).Decode(&addr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var deletedAddr bool
	for _, v := range md.Data {
		if v.Addr == addr.Addr {
			deletedAddr = true
		}
	}

	if deletedAddr {

		tmp := new(models.User)
		tmp.UID = md.UID
		tmp.Email = md.Email
		tmp.Pass = md.Pass
		tmp.Root = md.Root

		for _, v := range md.Data {
			if v.Addr == addr.Addr {
				continue
			}
			tmp.Data = append(tmp.Data, v)
		}

		// marshal for save to database
		var result []byte
		if result, err = tmp.MarshalBinary(); err != nil {
			log.Printf("error marshaling model: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}

		// Writing to bucket where cookie.value is uid in boltDB
		if err = core.Store.Write(tmp.UID, result); err != nil {
			log.Printf("error write user to database: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(core.GoodRequest)
	return
}
