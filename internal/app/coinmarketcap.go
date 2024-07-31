package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp,omitempty"`
		ErrorCode    int         `json:"error_code,omitempty"`
		ErrorMessage interface{} `json:"error_message,omitempty"`
		Elapsed      int         `json:"elapsed,omitempty"`
		CreditCount  int         `json:"credit_count,omitempty"`
		Notice       interface{} `json:"notice,omitempty"`
	} `json:"status,omitempty"`
	Data struct {
		ID          int       `json:"id,omitempty"`
		Symbol      string    `json:"symbol,omitempty"`
		Name        string    `json:"name,omitempty"`
		Amount      float64   `json:"amount,omitempty"`
		LastUpdated time.Time `json:"last_updated,omitempty"`
		Quote       struct {
			Usd struct {
				Price       float64   `json:"price,omitempty"`
				LastUpdated time.Time `json:"last_updated,omitempty"`
			} `json:"USD,omitempty"`
		} `json:"quote,omitempty"`
	} `json:"data,omitempty"`
}

func (core *Core) Course(pool string) (float64, error) {

	code := map[string]int{
		"etc":        1321,
		"solo-etc":   1321,
		"rvn":        2577,
		"solo-rvn":   2577,
		"xna":        27195,
		"solo-xna":   27195,
		"kas":        20396,
		"solo-kas":   20396,
		"erg":        1762,
		"solo-erg":   1762,
		"xmr":        328,
		"solo-xmr":   328,
		"btg":        2083,
		"solo-btg":   2083,
		"clore":      26497,
		"solo-clore": 26497,
		"ctxc":       2638,
		"solo-ctxc":  2638,
		"neox":       21045,
		"solo-neox":  21045,
		"ethw":       21296,
		"solo-ethw":  21296,
		"nexa":       23380,
		"solo-nexa":  23380,
		"kls":        29968,
		"solo-kls":   29968,
		"grin":       3709,
		"solo-grin":  3709,
		"beam":       3702,
		"solo-beam":  3702,
		"zec":        1437,
		"solo-zec":   1437,
		"firo":       1414,
		"solo-firo":  1414,
		"ckb":        4948,
		"solo-ckb":   4948,
		"flux":       3029,
		"solo-flux":  3029,
		"zen":        1698,
		"solo-zen":   1698,
		"ae":         1700,
		"solo-ae":    1700,
		"pyi":        29410,
		"solo-pyi":   29410,
		"bch":        1831,
		"solo-bch":   1831,
	}

	var (
		client http.Client
		res    *http.Response
		req    *http.Request
		err    error
	)

	url := fmt.Sprintf(
		"https://pro-api.coinmarketcap.com/v2/tools/price-conversion?id=%d&amount=%f", code[pool], 1.00)

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return 0.00, err
	}

	req.Header.Add("X-CMC_PRO_API_KEY", core.Config.CoinMarketCapApiKey)
	req.Header.Add("Accept", "*/*")

	if res, err = client.Do(req); err != nil {
		return 0.00, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("don't close body in coinmarketcap %v\n", err.Error())
		}
	}()

	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		return 0.00, err
	}

	var st *Response
	err = json.Unmarshal(body, &st)
	if err != nil {
		return 0.00, err
	}

	return st.Data.Quote.Usd.Price, nil
}
