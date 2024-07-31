package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/textproto"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/MarlikAlmighty/2miners/internal/models"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
)

func (core *Core) AccountRegisterHandler(rw http.ResponseWriter, r *http.Request) {

	form := new(models.FormLoginPassword)

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var (
		verifier = emailverifier.NewVerifier()
		ret      *emailverifier.Result
		err      error
	)

	// check email
	if ret, err = verifier.Verify(form.Login); err != nil {
		log.Printf("verify email address failed: %v\n", form.Login)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	// check email
	if !ret.Syntax.Valid {
		log.Printf("email address syntax is invalid: %v\n", form.Login)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	mp := make(map[string]models.User)
	var v models.User

	// Check if user registered
	if mp, err = core.Store.ReadAll("users"); err != nil {
		log.Printf("error read from database: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	for _, v = range mp {
		if form.Login == v.Email {
			log.Printf("error, there is already such a user: %v\n", form.Login)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	// Check if user wait validation
	if mp, err = core.Store.ReadAll("tokens"); err != nil {
		log.Printf("error read from database: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	var k string
	for k, v = range mp {
		if form.Login == v.Email {

			// Find user, send hash to email
			e := &email.Email{
				To:      []string{form.Login},
				From:    "Monitoring Bot <" + core.Config.SmtpUser + ">",
				Subject: "Registration in monitoring bot",
				HTML: []byte(
					"To <a href=https://" +
						core.Config.Domain + "/account/verify/" + k + ">confirm</a>" +
						" your registration, click on the link.",
				),
				Headers: textproto.MIMEHeader{},
			}

			if err = core.SendEmail(*e); err != nil {
				log.Printf("error send email: %v\n", form.Login)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(core.BadRequest)
				return
			}

			log.Printf("send double hash to email: %v\n", form.Login)

			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(core.BadRequest)
			return
		}
	}

	// New user

	var (
		user models.User
		b    []byte
	)

	uid := uuid.New()
	user.UID = uid.String()
	user.Email = form.Login
	user.Pass = core.StringToHash(form.Password)
	user.Root = false

	if b, err = user.MarshalBinary(); err != nil {
		log.Printf("error marshal model: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	h := sha256.New()
	h.Write([]byte(uid.String()))
	bs := h.Sum(nil)

	if err = core.Store.WriteTTL(fmt.Sprintf("%x", bs), b); err != nil {
		log.Printf("error write user to database: %v\n", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	// New user, send hash email
	e := &email.Email{
		To:      []string{form.Login},
		From:    "Monitoring Bot <" + core.Config.SmtpUser + ">",
		Subject: "Registration in monitoring bot",
		HTML: []byte(
			"To <a href=https://" +
				core.Config.Domain + "/account/verify/" + fmt.Sprintf("%x", bs) + ">confirm</a>" +
				" your registration, click on the link.",
		),
		Headers: textproto.MIMEHeader{},
	}

	if err = core.SendEmail(*e); err != nil {
		log.Printf("error send email: %v\n", form.Login)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(core.BadRequest)
		return
	}

	log.Printf("send hash to email: %v\n", form.Login)

	rw.WriteHeader(http.StatusOK)
	rw.Write(core.GoodRequest)
	return
}
