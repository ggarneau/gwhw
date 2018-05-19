package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggarneau/gateway/response"
)

//ProductHandler This could have more requester to other services if needed
type ProductHandler struct {
	ProductService Requester
}

//Product struct
type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//GetOne This is kinda stupid because i get the object, Decode it, and reencode it, but we might wanna do something with it, or return a different version of it
func (ph *ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	resp, err := ph.ProductService.Do(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		return
	}

	result := Product{}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Do stuff with Product if needed here. Or call other services.
	response.JSON(w, result)
}

//Post This is kinda stupid because i get the object, Decode it, and reencode it, but we might wanna do something with it, or return a different version of it
func (ph *ProductHandler) Post(w http.ResponseWriter, r *http.Request) {
	resp, err := ph.ProductService.Post(r.URL.Path, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
