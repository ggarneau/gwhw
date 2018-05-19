package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/ggarneau/gateway/response"
	"github.com/ggarneau/gateway/server/services"
	"github.com/gorilla/mux"
)

var port int

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&port, "port", 8082, "Port to listen to")
	f.Parse(arguments)
}

func main() {
	initFlags(flag.CommandLine, os.Args[1:])
	n := negroni.Classic()
	s := mux.NewRouter()
	s.HandleFunc("/company/{cid}/product/{id}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		response.JSON(rw, services.Product{ID: id, Name: "Product #" + vars["id"] + " Of Company #" + vars["cid"]})
	})
	s.HandleFunc("/company/{id}/product", func(rw http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		rw.WriteHeader(http.StatusCreated)
		rw.Write(b)
	})
	s.HandleFunc("/auth/validate", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	n.UseHandler(s)
	n.Run(fmt.Sprintf(":%d", port))
}
