package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"

	"github.com/ggarneau/gateway/client/config"
	"github.com/gorilla/mux"
)

var errConfig = errors.New("error in config file")

//NewRouter Start a configurable proxy server
func NewRouter(configFile string) (http.Handler, error) {
	s := mux.NewRouter()
	options, err := config.GetOptions(configFile)
	if err != nil {
		return nil, errConfig
	}
	for _, option := range options {
		if !option.Validate() {
			return nil, errConfig
		}
		h := NewClientHandler(option)
		s.Handle(option.Path, h).Methods(option.Methods...)
		log.Println("Listening on " + option.Path)
	}
	return s, nil
}

//Handler is our Proxy Handler
type Handler struct {
	*httputil.ReverseProxy
}

//NewClientHandler Gives back a client handler
func NewClientHandler(o config.Options) *Handler {

	//Director is used to modify the request URL
	director := func(req *http.Request) {
		//This is where we redirect the url
		URL, _ := url.Parse(getURL(o.RedirectURL, req))
		req.URL.Scheme = URL.Scheme
		req.URL.Host = URL.Host
		req.Header.Add("Authorization", "Bearer "+o.Token)
		//Dumb Proxy will keep the Path from the original
		if !o.DumbProxy {
			req.URL.Path = URL.Path
		}
	}

	//This sets the handler, which use the director above to proxy the request
	//The proxy has its own serveHTTP function which return the result
	h := &Handler{
		ReverseProxy: &httputil.ReverseProxy{
			Director: director,
			Transport: &Transport{
				option: o,
				RoundTripper: &http.Transport{
					//We want some settable timeouts depending of the path
					DialContext: (&net.Dialer{
						Timeout:   time.Duration(o.TimeoutInSecond) * time.Second,
						KeepAlive: time.Duration(o.KeepAliveInSecond) * time.Second,
						DualStack: true,
					}).DialContext,
				},
			},
		},
	}

	return h
}

//This function take variable sets in original url and set them on the redirect link
func getURL(url string, req *http.Request) string {
	vars := mux.Vars(req)
	for key, v := range vars {
		re := regexp.MustCompile(fmt.Sprintf("{%s}", key))
		url = re.ReplaceAllString(url, v)
	}
	return url
}

//Transport is used for custom made roundtripping
type Transport struct {
	http.RoundTripper
	option config.Options
}

//RoundTrip : This roundtripper does nothing but could potentially be modified if needed
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripper.RoundTrip(req)

}
