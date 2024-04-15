package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sf-mu/pkg/models/censors"

	"github.com/gorilla/mux"
)

//
type API struct {
	r  *mux.Router       
	db censors.Interface 
}

// creates censors api
func New(db censors.Interface) *API {
	api := API{
		r:  mux.NewRouter(),
		db: db,
	}
	api.endpoints()
	return &api
}

// return api router
func (api *API) Router() *mux.Router {
	return api.r
}

// creates api endpoints
func (api *API) endpoints() {
	api.r.HandleFunc("/comments/check", api.checkHandler).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments/stop", api.stopHandler).Methods(http.MethodPost, http.MethodOptions)
}

// responses with all black list rows
func (api *API) checkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	text := struct {
		Content string
	}{}
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blackList, err := api.db.AllList()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, stopWord := range blackList {
		matched, err := regexp.MatchString(stopWord.BanWord, text.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if matched {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// adds word to stop list
func (api *API) stopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c censors.BlackList
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = api.db.AddList(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.ResponseWriter.WriteHeader(w, http.StatusCreated)

}
