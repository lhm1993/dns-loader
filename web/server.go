package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
)

var currentPath string
var store *sessions.CookieStore

// NewServer function

func auth(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "cookie-name")
		if _, ok := session.Values["username"].(string); !ok {
			http.Redirect(w, req, "/login", 302)
			return
		}
		f(w, req)
	}

}

func index(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, "index", nil)
}

func startDNSTraffic(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	decoder := json.NewDecoder(req.Body)
	var config dnsloader.Configuration
	err := decoder.Decode(&config)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "decode config fail"})
	} else {
		// localTraffic
		go dnsloader.GenTrafficFromConfig(&config)
		r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
	}
}
func login(config *dnsloader.Configuration) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "cookie-name")
		if user, ok := session.Values["username"].(string); ok && user != "" {
			http.Redirect(w, req, "/index", 302)
			return
		}
		if req.Method == "POST" {
			user := req.FormValue("username")
			password := req.FormValue("password")
			if user == config.User && password == config.Password {
				session.Values["username"] = user
				session.Save(req, w)
				http.Redirect(w, req, "/index", 302)
			} else {
				http.Redirect(w, req, "/login", 401)
			}
		} else if req.Method == "GET" {
			r := render.New(render.Options{})
			r.HTML(w, http.StatusOK, "login", nil)
		}
	}
}

// NewServer function create the http api
func NewServer(config *dnsloader.Configuration) {
	key := []byte(config.AppSecrect)
	store = sessions.NewCookieStore(key)
	r := mux.NewRouter()
	r.HandleFunc("/", auth(index)).Methods("GET")
	r.HandleFunc("/login", login(config)).Methods("GET", "POST")
	r.HandleFunc("/start", auth(startDNSTraffic)).Methods("POST")
	log.Println("http server route init success")
	log.Printf("static file folder:%s\n", http.Dir("/web/assets"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	http.ListenAndServe(config.HTTPServer, r)
}
