package app

import (
	"log"
	"net/http"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

func (core *Core) AccountProfileHandler(rw http.ResponseWriter, r *http.Request) {

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
	} else {
		log.Printf("not found session, ip %v\n", r.RemoteAddr)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	md.UID = ""
	md.Pass = ""

	var addr models.Address
	for i, v := range md.Data {
		i++
		v.ID = int64(i)
		addr = append(addr, v)
	}

	md.Data = addr

	if b, err = md.MarshalBinary(); err != nil {
		log.Printf("error marshal model: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
	return
}
