package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

type StatsAddr struct {
	Pool string `json:"pool,omitempty"`
	Addr string `json:"addr,omitempty"`
}

func (core *Core) AddrGetStatsHandler(rw http.ResponseWriter, r *http.Request) {

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
	if len(ck.Value) > 0 {
		if b, err = core.Store.Read("users", value); err != nil {
			log.Printf("error read from database: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	// unmarshal data to model if we have user in database
	mp := new(models.User)
	if len(b) > 0 {
		if err = mp.UnmarshalBinary(b); err != nil {
			log.Printf("error unmarshal model: %v\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	statsAddr := StatsAddr{}
	if err = json.NewDecoder(r.Body).Decode(&statsAddr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var cur float64
	if cur, err = core.Course(statsAddr.Pool); err != nil {
		log.Printf("error from course %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var (
		bt       []byte
		splitter float64
	)

	switch statsAddr.Pool {

	case "grin", "solo-grin":
		bt, err = core.grin(statsAddr.Pool, statsAddr.Addr, cur)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(bt)
		return

	case "nexa", "solo-nexa":
		bt, err = core.nexa(statsAddr.Pool, statsAddr.Addr, cur)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(bt)
		return

	case "xmr", "solo-xmr":
		splitter = float64(1000000000000)

	case "rvn", "solo-rvn", "xna", "solo-xna", "btg", "solo-btg", "kas", "solo-kas", "clore", "solo-clore",
		"neox", "solo-neox", "beam", "solo-beam", "firo", "solo-firo", "zec", "solo-zec", "flux", "solo-flux",
		"pyi", "solo-pyi", "zen", "solo-zen", "ae", "solo-ae", "bch", "solo-bch", "kls", "solo-kls",
		"ckb", "solo-ckb":
		splitter = float64(100000000)

	default:
		splitter = float64(1000000000)
	}

	// http request with pool and addr for validation
	account := new(models.AccountReturnModel)
	if account, err = core.ReturnAccount(statsAddr.Pool, core.LowCase(statsAddr.Addr)); err != nil {
		log.Printf("error, wrong data: %v %v\n", statsAddr.Pool, statsAddr.Addr)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	for _, v := range mp.Data {
		if v.Pool == statsAddr.Pool && v.Addr == statsAddr.Addr {

			md := new(models.StatsAddr)

			// check if pool is solo
			if strings.HasPrefix(statsAddr.Pool, "solo-") {

				md.Balance = fmt.Sprintf("%.1f", float64(account.Stats.Balance)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Balance)/splitter)
				md.Immature = fmt.Sprintf("%.1f", float64(account.Stats.Immature)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Immature)/splitter)
				md.Payed = fmt.Sprintf("%.1f", float64(account.Stats.Paid)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Paid)/splitter)

				md.CurrentLuck = account.CurrentLuck
				md.LastBlockFound = account.Stats.BlocksFound

			} else {

				md.Balance = fmt.Sprintf("%.6f", float64(account.Stats.Balance)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Balance)/splitter)
				md.Immature = fmt.Sprintf("%.6f", float64(account.Stats.Immature)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Immature)/splitter)
				md.Payed = fmt.Sprintf("%.6f", float64(account.Stats.Paid)/splitter) + " " +
					fmt.Sprintf("%.2f $", cur*float64(account.Stats.Paid)/splitter)

			}

			switch statsAddr.Pool {
			case "btg", "solo-btg", "beam", "solo-beam", "zec", "solo-zec", "flux", "solo-flux", "zen", "solo-zen":
				md.HashRate = core.SolutionsHashRate(account.Hashrate)
				md.CurrentHashRate = core.SolutionsHashRate(account.CurrentHashrate)
			case "ctxc", "solo-ctxc", "grin", "solo-grin", "ae", "solo-ae":
				md.HashRate = core.GraphHashRate(account.Hashrate)
				md.CurrentHashRate = core.GraphHashRate(account.CurrentHashrate)
			default:
				md.HashRate = core.HashRate(account.Hashrate)
				md.CurrentHashRate = core.HashRate(account.CurrentHashrate)
			}

			for k, w := range account.Workers {
				var wrk models.Worker
				wrk.Name = k
				wrk.LastBeat = core.LastBeat(w.LastBeat)
				switch statsAddr.Pool {
				case "btg", "solo-btg", "beam", "solo-beam", "zec", "solo-zec", "flux", "solo-flux", "zen", "solo-zen":
					wrk.HR = core.SolutionsHashRate(w.Hr)
				case "ctxc", "solo-ctxc", "grin", "solo-grin", "ae", "solo-ae":
					wrk.HR = core.GraphHashRate(w.Hr)
				default:
					wrk.HR = core.HashRate(w.Hr)
				}
				md.Workers = append(md.Workers, &wrk)
			}

			for _, s := range account.Sumrewards {
				var sumRewards models.SumRewards
				if s.Inverval == 86400 || s.Inverval == 604800 {
					sumRewards.Name = s.Name
					if strings.HasPrefix(statsAddr.Pool, "solo-") {
						sumRewards.Reward = fmt.Sprintf("%.1f", float64(s.Reward)/splitter) + " " +
							fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)
					} else {
						sumRewards.Reward = fmt.Sprintf("%.6f", float64(s.Reward)/splitter) + " " +
							fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)
					}
					md.SumRewards = append(md.SumRewards, &sumRewards)
				}

				if s.Inverval == 2592000 {
					sumRewards.Name = s.Name
					if strings.HasPrefix(statsAddr.Pool, "solo-") {
						md.LastBlockFound = s.Numreward
						sumRewards.Reward = fmt.Sprintf("%.1f", float64(s.Reward)/splitter) + " " +
							fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)
					} else {
						sumRewards.Reward = fmt.Sprintf("%.6f", float64(s.Reward)/splitter) + " " +
							fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)

					}
					md.SumRewards = append(md.SumRewards, &sumRewards)
				}
			}

			var res []byte
			if res, err = md.MarshalBinary(); err != nil {
				log.Printf("error marshal model: %v\n", err.Error())
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return
			}

			rw.WriteHeader(http.StatusOK)
			rw.Write(res)
			return
		}
	}

	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(core.BadRequest)
	return
}

func (core *Core) nexa(pool, addr string, cur float64) ([]byte, error) {

	md := new(models.StatsAddr)
	account := new(models.AccountReturnModel)
	var err error

	account, err = core.ReturnAccount(pool, core.LowCase(addr))
	if err != nil {
		log.Printf("error, wrong data: %v %v %v\n", pool, addr, err.Error())
		return nil, err
	}

	splitter := float64(100)
	md.Balance = fmt.Sprintf("%.2f", float64(account.Stats.Balance)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Balance)/splitter)
	md.Immature = fmt.Sprintf("%.2f", float64(account.Stats.Immature)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Immature)/splitter)
	md.Payed = fmt.Sprintf("%.2f", float64(account.Stats.Paid)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Paid)/splitter)

	md.HashRate = core.HashRate(account.Hashrate)
	md.CurrentHashRate = core.HashRate(account.CurrentHashrate)

	if strings.HasPrefix(pool, "solo-") {
		md.LastBlockFound = account.Stats.BlocksFound
		md.CurrentLuck = account.CurrentLuck
	}

	for k, w := range account.Workers {
		var wrk models.Worker
		wrk.Name = k
		wrk.LastBeat = core.LastBeat(w.LastBeat)
		wrk.HR = core.HashRate(w.Hr)
		md.Workers = append(md.Workers, &wrk)
	}

	for _, s := range account.Sumrewards {
		var rw models.SumRewards
		if s.Inverval == 86400 || s.Inverval == 604800 || s.Inverval == 2592000 {
			rw.Name = s.Name
			rw.Reward = fmt.Sprintf("%.2f", float64(s.Reward)/splitter) + " " +
				fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)
			md.SumRewards = append(md.SumRewards, &rw)
		}
	}

	var res []byte
	if res, err = md.MarshalBinary(); err != nil {
		log.Printf("error marshal model: %v\n", err.Error())
		return nil, err
	}

	return res, nil
}

func (core *Core) grin(pool, addr string, cur float64) ([]byte, error) {

	md := new(models.StatsAddr)
	account := new(Account32)
	var err error

	account, err = core.ReturnAccount2(pool, core.LowCase(addr))
	if err != nil {
		log.Printf("error, wrong data: %v %v %v\n", pool, addr, err.Error())
		return nil, err
	}

	splitter := float64(1000000000)
	md.Balance = fmt.Sprintf("%.6f", float64(account.Stats.Balance)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Balance)/splitter)
	md.Immature = fmt.Sprintf("%.6f", float64(account.Stats.Immature)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Immature)/splitter)
	md.Payed = fmt.Sprintf("%.6f", float64(account.Stats.Paid)/splitter) + " " +
		fmt.Sprintf("%.2f $", cur*float64(account.Stats.Paid)/splitter)

	md.HashRate = fmt.Sprintf("%.2f G/ps", account.Hashrates.Num32)
	md.CurrentHashRate = fmt.Sprintf("%.2f G/ps", account.CurrentHashrates.Num32)

	if strings.HasPrefix(pool, "solo-") {
		md.LastBlockFound = int64(account.Stats.BlocksFound)
		md.CurrentLuck = account.CurrentLuck
	}

	for k, w := range account.Workers {
		var wrk models.Worker
		wrk.Name = k
		wrk.LastBeat = core.LastBeat(w.LastBeat)
		wrk.HR = fmt.Sprintf("%.2f G/ps", w.Hr)
		md.Workers = append(md.Workers, &wrk)
	}

	for _, s := range account.Sumrewards {
		var rw models.SumRewards
		if s.Inverval == 86400 || s.Inverval == 604800 || s.Inverval == 2592000 {
			rw.Name = s.Name
			rw.Reward = fmt.Sprintf("%.6f", float64(s.Reward)/splitter) + " " +
				fmt.Sprintf("%.2f $", cur*float64(s.Reward)/splitter)
			md.SumRewards = append(md.SumRewards, &rw)
		}
	}

	var res []byte
	if res, err = md.MarshalBinary(); err != nil {
		log.Printf("error marshal model: %v\n", err.Error())
		return nil, err
	}

	return res, nil
}
