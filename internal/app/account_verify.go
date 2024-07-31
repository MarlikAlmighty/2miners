package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"

	"github.com/MarlikAlmighty/2miners/internal/models"
)

type ViewData struct {
	Title   string
	Header  string
	Message string
}

var html = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/semantic-ui@2.5.0/dist/semantic.min.css">
    <style type="text/css">
        body>.grid {
            height: 100%;
        }
        .column {
            max-width: 450px;
        }
    </style>
    <title>{{ .Title}}</title>
</head>
<body>
    <div class="ui middle aligned center aligned grid">
        <div class="column">
            <div class="ui message positive">
                <div class="ui left aligned header">
                    {{ .Header}}
                </div>
                <div style="text-align: left;">
                    <p>{{ .Message}}</p>
                </div>
            </div>
        </div>
    </div>
</body>
<script>
    setTimeout(function () {
        window.location.href = "/";
    }, 5000); 
</script>
</html>`

func (core *Core) AccountVerifyHandler(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		fmt.Println("token is missing in parameters")
	}

	mp := make(map[string]models.User)
	var err error

	data := ViewData{
		Title:   "Status verification",
		Header:  "",
		Message: "",
	}

	tmpl := template.Must(template.New("data").Parse(html))

	if mp, err = core.Store.ReadAll("tokens"); err != nil {
		log.Printf("error reads users from database: %v\n", err.Error())
		data.Header = "Internal error"
		data.Message = http.StatusText(http.StatusBadRequest)
		if err = tmpl.Execute(rw, data); err != nil {
			log.Printf("error execute template %v", err.Error())
			return
		}
	}

	var (
		k string
		v models.User
	)

	for k, v = range mp {
		if token == k {
			var b []byte
			if b, err = v.MarshalBinary(); err != nil {
				data.Header = "Internal error"
				data.Message = http.StatusText(http.StatusBadRequest)
				log.Printf("error marshaling user: %v\n", err.Error())
				if err = tmpl.Execute(rw, data); err != nil {
					log.Printf("error execute template %v", err.Error())
					return
				}
			}

			if err = core.Store.Write(v.UID, b); err != nil {
				data.Header = "Internal error"
				data.Message = http.StatusText(http.StatusBadRequest)
				log.Printf("error write user to database: %v\n", err.Error())
				if err = tmpl.Execute(rw, data); err != nil {
					log.Printf("error execute template %v", err.Error())
					return
				}
			}

			if err = core.Store.Delete("tokens", token); err != nil {
				data.Header = "Internal error"
				data.Message = http.StatusText(http.StatusBadRequest)
				log.Printf("error delete token from database: %v\n", err.Error())
				if err = tmpl.Execute(rw, data); err != nil {
					log.Printf("error execute template %v", err.Error())
					return
				}
			}

			data.Header = "Success verification"
			data.Message = "You may <a href='/'>go to login page</a> or now we redirect over 5 seconds..."
			if err = tmpl.Execute(rw, data); err != nil {
				log.Printf("error execute template %v", err.Error())
				return
			}
		}
	}

	log.Printf("not found such user, ip: %v\n", r.RemoteAddr)

	data.Header = http.StatusText(http.StatusBadRequest)
	data.Message = "The link is invalid or out of date."
	if err = tmpl.Execute(rw, data); err != nil {
		log.Printf("error execute template %v", err.Error())
		return
	}
}
