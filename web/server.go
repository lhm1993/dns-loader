package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
)

var currentPath string
var store *sessions.CookieStore
var nodeManager *NodeManager

// NewServer function

func auth(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "dns-loader")
		if user, ok := session.Values["username"].(string); !ok || user == "" {
			http.Redirect(w, req, "/login", 302)
			return
		}
		f(w, req)
	}

}

func index(w http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		"iplist": nodeManager.IPList,
	}
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, "index", data)
}

func addNode(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	decoder := json.NewDecoder(req.Body)
	var ipinfo IPWithPort
	err := decoder.Decode(&ipinfo)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "decode data fail"})
		return
	}
	err = nodeManager.AddNode(ipinfo.IPAddress, ipinfo.Port)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func pingNode(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	decoder := json.NewDecoder(req.Body)
	var ipinfo IPWithPort
	err := decoder.Decode(&ipinfo)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	ip := fmt.Sprintf("%s:%d", ipinfo.IPAddress, ipinfo.Port)
	err = nodeManager.callPing(ip)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func deleteNode(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	decoder := json.NewDecoder(req.Body)
	var ipinfo IPWithPort
	err := decoder.Decode(&ipinfo)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	pending := fmt.Sprintf("%s:%d", ipinfo.IPAddress, ipinfo.Port)
	err = nodeManager.Remove(pending)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": err.Error()})
		return
	}
	r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
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
		err := config.Valid()
		if err != nil {
			log.Println(err)
			r.JSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "validate config fail"})
			return
		}
		go dnsloader.GenTrafficFromConfig(&config)
		go nodeManager.Call(Start, config)
		r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

func stopDNSTraffic(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	if stopStatus := dnsloader.GloablGenerator.Stop(); true != stopStatus {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "ServerFail"})
		return
	}
	go nodeManager.Call(Kill, nil)
	r.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func login(config *dnsloader.Configuration) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "dns-loader")
		if user, ok := session.Values["username"].(string); ok && user != "" {
			http.Redirect(w, req, "/", 302)
			return
		}
		if req.Method == "POST" {
			user := req.FormValue("username")
			password := req.FormValue("password")
			if user == config.User && password == config.Password {
				session.Values["username"] = user
				session.Save(req, w)
				http.Redirect(w, req, "/", 302)
			} else {
				http.Redirect(w, req, "/login", 401)
			}
		} else if req.Method == "GET" {
			r := render.New(render.Options{})
			r.HTML(w, http.StatusOK, "login", nil)
		}
	}
}
func logout(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "dns-loader")
	if user, ok := session.Values["username"].(string); ok && user != "" {
		session.Values["username"] = ""
		session.Save(req, w)
		http.Redirect(w, req, "/login", 302)
		return
	}
}

// NewServer function create the http api
func NewServer(config *dnsloader.Configuration) {
	key := []byte(config.AppSecrect)
	nodeManager = NewNodeManager(config)
	store = sessions.NewCookieStore(key)
	r := mux.NewRouter()
	r.HandleFunc("/", auth(index)).Methods("GET")
	r.HandleFunc("/logout", logout).Methods("POST", "GET")
	r.HandleFunc("/login", login(config)).Methods("GET", "POST")
	r.HandleFunc("/nodes", auth(addNode)).Methods("POST")
	r.HandleFunc("/nodes", auth(deleteNode)).Methods("DELETE")
	r.HandleFunc("/ping", auth(pingNode)).Methods("POST")
	r.HandleFunc("/start", auth(startDNSTraffic)).Methods("POST")
	r.HandleFunc("/stop", auth(stopDNSTraffic)).Methods("GET")
	log.Println("http server route init success")
	log.Printf("static file folder:%s\n", http.Dir("/web/assets"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	http.ListenAndServe(config.HTTPServer, http.TimeoutHandler(r, time.Second*10, "timeout"))
}
