package app

import (
	"encoding/json"
	"errors"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"io"
	"log"
	"net/http"
)

type Account32 struct {
	Two4Hnumreward int   `json:"24hnumreward"`
	Two4Hreward    int64 `json:"24hreward"`
	Charts         struct {
		Num32 []struct {
			MinerHash      float64 `json:"minerHash"`
			MinerLargeHash float64 `json:"minerLargeHash"`
			TimeFormat     string  `json:"timeFormat"`
			WorkerOnline   string  `json:"workerOnline"`
			X              int     `json:"x"`
		} `json:"32"`
	} `json:"charts"`
	CurrentHashrates struct {
		Num32 float32 `json:"32"`
	} `json:"currentHashrates"`
	CurrentLuck string `json:"currentLuck"`
	Hashrates   struct {
		Num32 float64 `json:"32"`
	} `json:"hashrates"`
	PageSize int `json:"pageSize"`
	Payments []struct {
		Amount         int64  `json:"amount"`
		Timestamp      int    `json:"timestamp"`
		Tx             string `json:"tx"`
		TxFee          int    `json:"txFee"`
		TxKernelExcess string `json:"txKernelExcess"`
	} `json:"payments"`
	PaymentsTotal int `json:"paymentsTotal"`
	Rewards       []struct {
		Blockheight int     `json:"blockheight"`
		Timestamp   int     `json:"timestamp"`
		Blockhash   string  `json:"blockhash"`
		Reward      int     `json:"reward"`
		Percent     float64 `json:"percent"`
		Immature    bool    `json:"immature"`
		Orphan      bool    `json:"orphan"`
	} `json:"rewards"`
	RoundShares int `json:"roundShares"`
	Stats       struct {
		Balance     int64 `json:"balance"`
		BlocksFound int   `json:"blocksFound"`
		Immature    int64 `json:"immature"`
		LastShare   int   `json:"lastShare"`
		Paid        int64 `json:"paid"`
		Pending     int   `json:"pending"`
	} `json:"stats"`
	Sumrewards []struct {
		Inverval  int    `json:"inverval"`
		Reward    int    `json:"reward"`
		Numreward int    `json:"numreward"`
		Name      string `json:"name"`
		Offset    int    `json:"offset"`
	} `json:"sumrewards"`
	UpdatedAt      int64                                `json:"updatedAt"`
	Workers        map[string]*models.WorkerGroupModel2 `json:"workers,omitempty"`
	WorkersOffline int                                  `json:"workersOffline"`
	WorkersOnline  int                                  `json:"workersOnline"`
	WorkersTotal   int                                  `json:"workersTotal"`
}

func (core *Core) ReturnAccount2(pool, addr string) (*Account32, error) {

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
	account := new(Account32)
	if err = json.Unmarshal(b, account); err != nil {
		return nil, err
	}

	return account, nil
}
