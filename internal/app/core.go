package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/securecookie"

	"github.com/jordan-wright/email"

	"github.com/MarlikAlmighty/2miners/internal/models"

	"github.com/MarlikAlmighty/2miners/internal/store"

	"github.com/MarlikAlmighty/2miners/internal/config"
)

// Core main app
type Core struct {
	Server      *http.Server               `server:"-"`
	Store       Store                      `store:"-"`
	Config      *config.Configuration      `config:"-"`
	Secure      *securecookie.SecureCookie `secure:"-"`
	BadRequest  []byte
	GoodRequest []byte
}

// New core app
func New(c *config.Configuration, s *store.Wrapper, p *securecookie.SecureCookie) *Core {

	bReq := &models.ExceptionModel{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}

	b, err := bReq.MarshalBinary()
	if err != nil {
		log.Fatalf("error marshall: %v\n", err.Error())
	}

	gReq := &models.ExceptionModel{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}

	var g []byte
	if g, err = gReq.MarshalBinary(); err != nil {
		log.Fatalf("error marshall: %v\n", err.Error())
	}

	return &Core{
		Config:      c,
		Store:       s,
		Secure:      p,
		BadRequest:  b,
		GoodRequest: g,
		Server:      &http.Server{},
	}
}

// Run core app
func (core *Core) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// catch signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// stop core server
	go func() {
		for {
			select {
			case <-shutdown:
				cancel()
			}
		}
	}()

	go core.StartServer()
	go core.Clean(ctx)

	var (
		mp  = make(map[string]models.User)
		md  models.User
		wg  sync.WaitGroup
		err error
	)

	// for number of requests per second
	ticker := time.NewTicker(time.Duration(core.Config.RequestOverTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err = core.Store.Close(); err != nil {
				log.Printf("error close database: %v\n", err)
			}
			if err = core.Server.Shutdown(context.Background()); err != nil {
				log.Printf("server shutdown error: %v\n", err.Error())
			}
			return
		case <-ticker.C:
			if mp, err = core.Store.ReadAll("users"); err != nil {
				log.Fatalf("error read database: %v\n", err)
			}
			for _, md = range mp {
				if len(md.Data) > 0 {
					wg.Add(1)
					go func() {
						defer wg.Done()
						core.LifeTime(md)
						core.Block(md)
					}()
					wg.Wait()
				}
			}
		}
	}
}

// LifeTime monitors the life cycle of rigs
func (core *Core) LifeTime(user models.User) models.User {

	var (
		tmp        models.Address
		client     http.Client
		res        *http.Response
		req        *http.Request
		t, n, tNow int64
		err        error
	)

	for _, data := range user.Data {

		var del bool

		if req, err = http.NewRequest("GET",
			"https://"+data.Pool+".2miners.com/api/accounts/"+core.LowCase(data.Addr), nil); err != nil {
			log.Printf("[WORKER]: error new request: %s\n", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json")
		if res, err = client.Do(req); err != nil {
			log.Printf("[WORKER]: error make request %v\n", err.Error())
			continue
		}

		if res.StatusCode != 200 {
			log.Printf("[WORKER]: error get api, status %v\n", res.StatusCode)
			continue
		}

		var b []byte
		if b, err = io.ReadAll(res.Body); err != nil {
			log.Printf("[WORKER]: could not read response body: %s\n", err.Error())
			continue
		}

		if err = res.Body.Close(); err != nil {
			log.Printf("[WORKER]: could not close body: %s\n", err.Error())
			continue
		}

		// for response api
		account := new(models.AccountReturnModel)
		if err = json.Unmarshal(b, account); err != nil {
			log.Printf("[WORKER]: error unmarshal response body: %s\n", err.Error())
			continue
		}

		// we have notifies in minutes convert it to seconds and minus timeout
		n = (data.Notify * 60) - int64(core.Config.RequestOverTime)
		tNow = time.Now().Unix()

		for i, w := range account.Workers {

			// time now minus timestamp from api
			t = tNow - w.LastBeat

			// send alert to email
			if t > n {

				del = true

				if data.MonitorAddr {

					e := &email.Email{
						To:      []string{user.Email},
						From:    "Monitoring Bot <" + core.Config.SmtpUser + ">",
						Subject: "Alert from monitoring bot",
						HTML: []byte("Attention! Your worker: " + i + " offline, on pool: " + data.Pool +
							". We are suspending monitoring of your address. You can turn it back on in the settings."),
						Headers: textproto.MIMEHeader{},
					}

					if err = core.SendEmail(*e); err != nil {
						log.Printf("[WORKER]: could not send email %s\n", err.Error())
						break
					}
				}
			}
		}

		if del {
			data.MonitorAddr = false
			tmp = append(tmp, data)
		} else {
			tmp = append(tmp, data)
		}
	}

	user.Data = tmp
	if err = core.WriteUserToDataBase(user); err != nil {
		log.Printf("[CORE]: error: %v\n", err.Error())
	}

	return user
}

// Block monitors user blocks
func (core *Core) Block(user models.User) {

	type Blocks struct {
		FindersBlock []struct {
			Finders   string      `json:"finders,omitempty"`
			Hash      string      `json:"hash,omitempty"`
			Height    interface{} `json:"height,omitempty"`
			Timestamp interface{} `json:"timestamp,omitempty"`
		} `json:"findersBlock,omitempty"`
	}

	var (
		blocks Blocks
		client http.Client
		res    *http.Response
		req    *http.Request
		err    error
	)

	var write = false
	tmp := user

	for _, data := range tmp.Data {

		if data.MonitorBlock {

			if req, err = http.NewRequest("GET",
				"https://"+data.Pool+".2miners.com/api/finders", nil); err != nil {
				log.Printf("[BLOCK]: error new request: %s\n", err.Error())
				continue
			}

			req.Header.Set("Accept", "application/json")
			if res, err = client.Do(req); err != nil {
				log.Printf("[BLOCK]: error make request %v\n", err.Error())
				continue
			}

			if res.StatusCode != 200 {
				log.Printf("[BLOCK]: error get api, status %v\n", res.StatusCode)
				continue
			}

			var b []byte
			if b, err = io.ReadAll(res.Body); err != nil {
				log.Printf("[BLOCK]: could not read response body: %s\n", err.Error())
				continue
			}

			if err = res.Body.Close(); err != nil {
				log.Printf("[BLOCK]: could not close body: %s\n", err.Error())
				continue
			}

			if err = json.Unmarshal(b, &blocks); err != nil {
				log.Printf("[BLOCK]: could not unmarshal to model: %s\n", err.Error())
				continue
			}

			var send = false

			for i := range blocks.FindersBlock {

				switch typ := blocks.FindersBlock[i].Height.(type) {

				case string:
					var tmpBlock int64
					if tmpBlock, err = strconv.ParseInt(typ, 10, 64); err != nil {
						log.Printf("[BLOCK]: could not atoi to int: %s\n", err.Error())
						break
					}
					if tmpBlock > data.Block &&
						core.LowCase(blocks.FindersBlock[i].Finders) == core.LowCase(data.Addr) {
						data.Block = tmpBlock
						send = true
					}
				case int64:
					if typ > data.Block &&
						core.LowCase(blocks.FindersBlock[i].Finders) == core.LowCase(data.Addr) {
						data.Block = typ
						send = true
					}
				case float64:
					if int64(typ) > data.Block &&
						core.LowCase(blocks.FindersBlock[i].Finders) == core.LowCase(data.Addr) {
						data.Block = int64(typ)
						send = true
					}
				default:
					log.Printf("[BLOCK]: unknown type %T %v\n", typ, typ)
					break
				}

				if send {

					// write to database
					write = true

					e := &email.Email{
						To:      []string{user.Email},
						From:    "Monitoring Bot <" + core.Config.SmtpUser + ">",
						Subject: "Alert from monitoring bot",
						HTML: []byte("Attention! Block found at this address: " +
							data.Addr + " on this pool: " +
							data.Pool + "."),
						Headers: textproto.MIMEHeader{},
					}

					if err = core.SendEmail(*e); err != nil {
						log.Printf("[BLOCK]: could not send email %s\n", err.Error())
						break
					}
				}
				break
			}

			if write {
				if err = core.WriteUserToDataBase(tmp); err != nil {
					log.Printf("[CORE]: error: %v\n", err.Error())
				}
			}
		}
	}
}
