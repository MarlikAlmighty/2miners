package main

import (
	"github.com/MarlikAlmighty/2miners/internal/app"
	"github.com/MarlikAlmighty/2miners/internal/config"
	"github.com/MarlikAlmighty/2miners/internal/store"
	"github.com/gorilla/securecookie"
	"log"
)

func main() {

	// get the configuration for the app through ENV
	cnf := config.New()
	if err := cnf.GetEnv(); err != nil {
		log.Fatalf("get environment keys error: %v\n", err)
	}

	// connect to database
	r, err := store.New()
	if err != nil {
		log.Fatalf("open database error: %v\n", err)
	}

	// get hash for cookie
	s := securecookie.New([]byte(cnf.CookieHashKey), []byte(cnf.CookieBlockKey))

	// new app, start app
	core := app.New(cnf, r, s)
	core.Run()
}
