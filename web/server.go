package web

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
)

func initRoutes(router *mux.Router, config *dnsloader.Configuration) {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Panicf("Server path not valid:%s", err.Error())
	}
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		session := sessions.GetSession(r)
		if r.Method == "POST" {
			user := r.FormValue("username")
			password := r.FormValue("password")
			if user == config.User && password == config.Password {
				session.Set("username", user)
				http.Redirect(w, r, "/", 302)
			} else {
				http.Redirect(w, r, "/login", 200)
			}
		} else if r.Method == "GET" {
			r := render.New(render.Options{})
			r.HTML(w, http.StatusOK, "login", nil)
		}
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(currentPath + "/web/assets")))
}

// NewServer function
func NewServer(config *dnsloader.Configuration) *negroni.Negroni {
	n := negroni.Classic()
	router := mux.NewRouter()
	initRoutes(router, config)
	n.UseHandler(router)
	return n
}
