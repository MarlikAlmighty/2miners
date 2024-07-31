package app

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/jordan-wright/email"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

type (
	App interface {
		SolutionsHashRate(value float32) string
		GraphHashRate(value float32) string
		ClearAddr(str string) string
		Clean(ctx context.Context)
		WriteUserToDataBase(user models.User) error
		LastBeat(value int64) string
		HashRate(value float32) string
		LowCase(str string) string
		ReturnAccount(pool, addr string) (*models.AccountReturnModel, error)
		SendEmail(e email.Email) error
		StringToHash(s string) string
	}
	Store interface {
		Write(uid string, value []byte) error
		WriteTTL(uid string, value []byte) error
		Read(bucket, uid string) ([]byte, error)
		ReadAll(bucket string) (map[string]models.User, error)
		Delete(bucket, uid string) error
		GetExpired(maxAge time.Duration) ([][]byte, error)
		Sweep(keys [][]byte)
		Close() error
	}
	Config interface {
	}
)

func (core *Core) SolutionsHashRate(value float32) string {
	if value < 1000 {
		return fmt.Sprintf("%.2f\tS/s", value)
	}
	var arr = []string{"KS/s", "MS/s", "GS/s", "TS/s", "PS/s", "ES/s", "ZS/s"}
	i := 0
	for value > 1000 {
		value = value / 1000
		i++
	}
	return fmt.Sprintf("%.2f", value) + " " + arr[i-1]
}

func (core *Core) GraphHashRate(value float32) string {
	if value < 1000 {
		return fmt.Sprintf("%.2f\tGp/s", value)
	}
	var arr = []string{"KGp/s", "MGp/s", "GGp/s", "TGp/s", "PGp/s", "EGp/s", "ZGp/s"}
	i := 0
	for value > 1000 {
		value = value / 1000
		i++
	}
	return fmt.Sprintf("%.2f", value) + " " + arr[i-1]
}

func (core *Core) HashRate(value float32) string {
	if value < 1000 {
		return fmt.Sprintf("%.2f\tHH/s", value)
	}
	var arr = []string{"KH/s", "MH/s", "GH/s", "TH/s", "PH/s", "EH/s", "ZH/s"}
	i := 0
	for value > 1000 {
		value = value / 1000
		i++
	}
	return fmt.Sprintf("%.2f", value) + " " + arr[i-1]
}

func (core *Core) ClearAddr(str string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9:_ ]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func (core *Core) Clean(ctx context.Context) {

	var (
		keys [][]byte
		err  error
	)

	clearTicker := time.NewTicker(1 * time.Hour)
	defer clearTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-clearTicker.C:
			keys, err = core.Store.GetExpired(1 * time.Hour)
			if err != nil {
				log.Printf("error get expired tokens %v\n", err.Error())
				return
			}
			if len(keys) > 0 {
				core.Store.Sweep(keys)
			}
		}
	}
}

func (core *Core) WriteUserToDataBase(user models.User) error {
	var (
		result []byte
		err    error
	)
	if result, err = user.MarshalBinary(); err != nil {
		return err
	}
	if err = core.Store.Write(user.UID, result); err != nil {
		return err
	}
	return nil
}

func (core *Core) LastBeat(value int64) string {
	tNow := time.Now().Unix()
	tNow = tNow - value
	if tNow >= 3600 {
		return fmt.Sprintf("%v hours", tNow/3600)
	} else if tNow >= 60 {
		return fmt.Sprintf("%v min", tNow/60)
	}
	return fmt.Sprintf("%v sec", tNow)
}

func (core *Core) LowCase(str string) string {
	if strings.HasPrefix(str, "0x") {
		return strings.ToLower(str)
	}
	return str
}

func (core *Core) SendEmail(e email.Email) error {

	a := smtp.PlainAuth(
		"",
		core.Config.SmtpUser,
		core.Config.SmtpPassword,
		core.Config.SmtpHost)

	if err := e.Send(core.Config.SmtpHost+":"+core.Config.SmtpPort, a); err != nil {
		return err
	}
	return nil
}

func (core *Core) StringToHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
