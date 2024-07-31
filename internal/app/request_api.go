package app

import (
	"encoding/json"
	"errors"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"io"
	"log"
	"net/http"
)

func (core *Core) ReturnAccount(pool, addr string) (*models.AccountReturnModel, error) {

	var (
		client http.Client
		res    *http.Response
		req    *http.Request
		err    error
	)

	if req, err = http.NewRequest("GET",
		"https://"+pool+".2miners.com/api/accounts/"+addr, nil); err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("bad request")
	}

	var b []byte
	if b, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("don't close body in ReturnAccount %v\n", err.Error())
		}
	}()

	// for response api
	account := new(models.AccountReturnModel)
	if err = json.Unmarshal(b, account); err != nil {
		return nil, err
	}

	return account, nil
}
