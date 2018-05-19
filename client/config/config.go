package config

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
)

type Options struct {
	Path              string   `json:"path"`
	DumbProxy         bool     `json:"dumb_proxy"`
	RedirectURL       string   `json:"redirect"`
	TimeoutInSecond   int64    `json:"timeout"`
	KeepAliveInSecond int64    `json:"keepAlive"`
	Methods           []string `json:"methods"`
	Token             string   `json:"token"`
}

func GetOptions(config string) ([]Options, error) {
	var o []Options
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(raw, &o)
	return o, err
}

func (o Options) Validate() bool {
	if o.TimeoutInSecond == 0 {
		o.TimeoutInSecond = 30
	}
	if o.KeepAliveInSecond == 0 {
		o.KeepAliveInSecond = 30
	}
	if o.Path == "" || o.RedirectURL == "" || len(o.Methods) == 0 {
		return false
	}
	_, err := url.Parse(o.RedirectURL)
	return err == nil
}
