package main

import (
	"flag"
	"os"

	"github.com/ggarneau/gateway/server"
)

var port int
var jwtsecret string
var service string

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&port, "port", 8081, "Port to listen to")
	f.StringVar(&jwtsecret, "jwtsecret", "secret", "JWTSecret")
	f.StringVar(&service, "service", "http://localhost:8082", "Microservice")
	f.Parse(arguments)
}

func main() {
	initFlags(flag.CommandLine, os.Args[1:])
	s, err := server.NewServer(port, jwtsecret, service)
	if err != nil {
		panic(err)
	}
	s.Start()
}
