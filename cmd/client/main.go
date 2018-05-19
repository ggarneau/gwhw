package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/ggarneau/gateway/client"
)

var port int
var config string

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&port, "port", 8080, "Port to listen to")
	f.StringVar(&config, "config", "conf.json", "Config file location")
	f.Parse(arguments)
}

func main() {
	initFlags(flag.CommandLine, os.Args[1:])
	r, err := client.NewRouter(config)
	n := negroni.Classic()
	if err != nil {
		log.Panic(err)
	}
	port := fmt.Sprintf(":%d", port)
	n.UseHandler(r)
	n.Run(port)
}
