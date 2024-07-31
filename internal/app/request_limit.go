package app

import (
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

// LimitRequest limit request from handlers
func LimitRequest(next http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetMethods([]string{"GET", "POST", "UPDATE", "DELETE"})
	middle := func(w http.ResponseWriter, r *http.Request) {
		httpError := tollbooth.LimitByRequest(lmt, w, r)
		if httpError != nil {
			lmt.ExecOnLimitReached(w, r)
			w.Header().Add("Content-Type", lmt.GetMessageContentType())
			w.WriteHeader(httpError.StatusCode)
			if _, err := w.Write([]byte(httpError.Message)); err != nil {
				log.Printf("error write message to user: %v\n", r.RemoteAddr)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(middle)
}
