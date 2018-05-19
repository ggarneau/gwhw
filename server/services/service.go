package services

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

//SimpleService We could create more complex Requesters for more complicated services (async etc)
type SimpleService struct {
	Host        string
	ContentType string
	*http.Client
}

//Requester http tool interface
type Requester interface {
	Do(*http.Request) (*http.Response, error)
	Get(string) (*http.Response, error)
	Post(string, io.Reader) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
	Head(string) (*http.Response, error)
}

//Do request
func (p *SimpleService) Do(req *http.Request) (*http.Response, error) {
	//home made proxy
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	URL, _ := url.Parse(p.Host)
	url := fmt.Sprintf("%s://%s%s", URL.Scheme, URL.Host, req.RequestURI)
	proxyReq, _ := http.NewRequest(req.Method, url, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}
	return p.Client.Do(proxyReq)
}

//Get get method
func (p *SimpleService) Get(path string) (*http.Response, error) {
	//this is extremely simplefied but we could get multiple GETS and some Post here to multiple services
	return p.Client.Get(p.Host + path)
}

//Post post method
func (p *SimpleService) Post(path string, body io.Reader) (*http.Response, error) {
	return p.Client.Post(p.Host+path, p.ContentType, body)
}

//PostForm Post a simple form
func (p *SimpleService) PostForm(path string, data url.Values) (*http.Response, error) {
	return p.Client.PostForm(p.Host+path, data)
}

//Head get header
func (p *SimpleService) Head(path string) (*http.Response, error) {
	return p.Client.Head(p.Host + path)
}
