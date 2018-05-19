package negronimux

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func NewServer(conf Config) (*Server, error) {
	Router := mux.NewRouter().StrictSlash(true)

	if conf.NotFoundHandler != nil {
		Router.NotFoundHandler = conf.NotFoundHandler
	}

	if conf.Negroni == nil {
		conf.Negroni = negroni.Classic()
	}
	server := Server{
		&conf,
		Router,
	}
	return &server, nil
}

type Server struct {
	*Config
	Router *mux.Router
}

func (s *Server) RouteGroup(path string, routing func(*Group)) {
	g := &Group{
		PathPrefix: s.PathPrefix + path,
		Router:     mux.NewRouter(),
		Negroni:    negroni.New(),
	}
	routing(g)
	g.Negroni.Use(negroni.Wrap(g.Router))
	s.Router.PathPrefix(s.PathPrefix + path).Handler(g.Negroni)
}

func (s *Server) Start() {
	s.Negroni.Use(negroni.Wrap(s.Router))
	s.Negroni.Run(fmt.Sprintf(":%d", s.Port))
}

func (s *Server) Middleware(ms ...negroni.Handler) {
	for _, m := range ms {
		s.Negroni.Use(m)
	}
}

func (s *Server) Route(r Route) {
	s.Router.Path(s.PathPrefix + r.Path).Methods(r.Method).Handler(r.Handler)

}
