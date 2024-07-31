package app

import (
	"encoding/json"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"log"
	"net/http"
)

func (core *Core) AddrUpdateHandler(rw http.ResponseWriter, r *http.Request) {

	// Check cookie
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

	addrModel := new(models.Addr)
	if err = json.NewDecoder(r.Body).Decode(&addrModel); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	// here limits for saving addresses
	if addrModel.Notify <= 0 || addrModel.Notify > 60 {
		log.Printf("error, wrong data: %v\n", r.Form.Get("Notify"))
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

	for _, v := range md.Data {
		if v.Addr == addrModel.Addr {
			if addrModel.Notify > 0 {
				v.Notify = addrModel.Notify
			}
			v.MonitorAddr = addrModel.MonitorAddr
			v.MonitorBlock = addrModel.MonitorBlock

			// marshal for save to database
			var result []byte
			if result, err = md.MarshalBinary(); err != nil {
				log.Printf("error marshal model: %v\n", err.Error())
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return
			}

			// Writing to database
			if err = core.Store.Write(md.UID, result); err != nil {
				log.Printf("error write user to database: %v\n", err.Error())
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return
			}

			rw.WriteHeader(http.StatusOK)
			rw.Write(core.GoodRequest)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(core.GoodRequest)
	return
}
