package server

import (
	"net/http"
	"time"

	"github.com/ggarneau/gateway/server/auth"
	"github.com/ggarneau/gateway/server/negronimux"
	"github.com/ggarneau/gateway/server/services"
)

func NewServer(port int, jwtsecret string, host string) (*negronimux.Server, error) {
	server, err := negronimux.NewServer(negronimux.Config{
		Port: port,
	})
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Timeout: 3 * time.Second,
	}
	s := &services.SimpleService{
		Host:   host,
		Client: c,
	}
	p := services.ProductHandler{
		ProductService: s,
	}
	server.RouteGroup("/company/{id}", func(rg *negronimux.Group) {
		rg.Middleware(auth.NewAuthMiddleware(jwtsecret))
		rg.Middleware(auth.AuthorizationMiddleware{Client: c, AuthorizationServiceHost: host + "/auth/validate"})
		rg.Route(negronimux.Route{
			Method:  "POST",
			Path:    "/product",
			Handler: http.HandlerFunc(p.Post),
		})
		rg.Route(negronimux.Route{
			Method:  "GET",
			Path:    "/product/{id}",
			Handler: http.HandlerFunc(p.GetOne),
		})
	})
	return server, nil
}
