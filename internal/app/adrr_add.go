package app

import (
	"encoding/json"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"log"
	"net/http"
)

func (core *Core) AdrrAddHandler(rw http.ResponseWriter, r *http.Request) {

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

	addrModel := new(models.Addr)
	if err = json.NewDecoder(r.Body).Decode(&addrModel); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	// unmarshal data to model if we have user in database
	md := new(models.User)
	if len(b) > 0 {

		// we have user in database
		if err = md.UnmarshalBinary(b); err != nil {
			log.Printf("error unmarshal model: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}

		// check pool and addr in database
		for k, v := range md.Data {
			if addrModel.Pool == v.Pool && addrModel.Addr == v.Addr {
				log.Printf("address is already in the database: %v %v\n",
					addrModel.Pool, addrModel.Addr)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return

				// check limit here
			} else if k == core.Config.MaxAddr {
				log.Printf("error, limit address: %v\n", core.Config.MaxAddr)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return
			}
		}

		Addr := core.LowCase(core.ClearAddr(addrModel.Addr))

		// http request with pool and addr for validation
		if _, err = core.ReturnAccount(addrModel.Pool, Addr); err != nil {

			log.Printf("error, wrong data: %v %v user: %v\n",
				addrModel.Pool,
				Addr,
				md.Email)

			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}

		addr := new(models.Addr)
		addr.Addr = Addr
		addr.Pool = core.LowCase(addrModel.Pool)
		addr.MonitorBlock = addrModel.MonitorBlock
		addr.MonitorAddr = addrModel.MonitorAddr
		addr.Notify = addrModel.Notify
		addr.Block = 0

		md.Data = append(md.Data, addr)

		// marshal for save to database
		var result []byte
		if result, err = md.MarshalBinary(); err != nil {
			log.Printf("error marshaling user: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}

		// Writing user to database where cookie value is uid user in boltDB
		if err = core.Store.Write(md.UID, result); err != nil {
			log.Printf("error write user to database: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}

		// success, return status 200
		rw.WriteHeader(http.StatusOK)
		rw.Write(core.GoodRequest)
		return
	}

	log.Printf("user not found in database: %v\n", r.RemoteAddr)
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(core.BadRequest)
	return
}
