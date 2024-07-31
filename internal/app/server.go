package app

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (core *Core) StartServer() {

	fs := http.FileServer(http.Dir("./build"))
	r := mux.NewRouter()

	r.PathPrefix("/static/").HandlerFunc(fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/asset-manifest.json", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/favicon.ico", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/logo192.png", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/logo512.png", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/manifest.json", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/robots.txt", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/service-worker.js", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/service-worker.js.map", fs.ServeHTTP).Methods("GET")
	r.HandleFunc("/", fs.ServeHTTP).Methods("GET")

	r.HandleFunc("/account", core.AccountProfileHandler).Methods("GET")
	r.HandleFunc("/account", core.AccountDeleteHandler).Methods("DELETE")
	r.HandleFunc("/account/exit", core.AccountExitHandler).Methods("GET")
	r.HandleFunc("/account/login", core.AccountLoginHandler).Methods("POST")
	r.HandleFunc("/account/register", core.AccountRegisterHandler).Methods("POST")
	r.HandleFunc("/account/verify/{token}", core.AccountVerifyHandler).Methods("GET")

	r.HandleFunc("/address", core.AddrUpdateHandler).Methods("PUT")
	r.HandleFunc("/address", core.AdrrAddHandler).Methods("POST")
	r.HandleFunc("/address", core.AddrDeleteHandler).Methods("DELETE")
	r.HandleFunc("/address/state", core.AddrGetStateHandler).Methods("POST")
	r.HandleFunc("/address/stats", core.AddrGetStatsHandler).Methods("POST")

	r.HandleFunc("/users", core.UsersGetHandler).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(core.NotFound)

	core.Server.Addr = "0.0.0.0:8080"
	core.Server.WriteTimeout = 15 * time.Second
	core.Server.ReadTimeout = 15 * time.Second
	core.Server.IdleTimeout = 60 * time.Second
	core.Server.Handler = LimitRequest(r)

	log.Printf("Serve on %v\n", core.Server.Addr)
	err := core.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Serve error: %v\n", err.Error())
	}
}
