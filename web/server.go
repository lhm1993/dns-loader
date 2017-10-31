package web

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func initRoutes(mx *mux.Router) {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Panicf("Server path not valid:%s", err.Error())
	}
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(currentPath + "/web/assets")))
}

// NewServer function
func NewServer() *negroni.Negroni {
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx)
	n.UseHandler(mx)
	return n
}
