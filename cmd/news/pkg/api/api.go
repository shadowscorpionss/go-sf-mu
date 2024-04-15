package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sf-mu/pkg/models/news"

	"strconv"

	"github.com/gorilla/mux"
)


type API struct {
	r  *mux.Router    //api router
	db news.Interface 
}

const limit = 10

// creates new api structure
func New(db news.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// returns api router
func (api *API) Router() *mux.Router {
	return api.r
}

// creates api endpoints
func (api *API) endpoints() {

	// 
	api.r.HandleFunc("/news", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/latest", api.newsLatestHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/search", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
}

// postsHandler responses with all posts
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	sParam := r.URL.Query().Get("s")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, pagination, err := api.db.PostSearchILIKE(sParam, limit, (page-1)*limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Posts      []news.Post
		Pagination news.Pagination
	}{
		Posts:      posts,
		Pagination: pagination,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// responses with page of posts
// query contains page number
func (api *API) newsLatestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.db.Posts(limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// responses with detailed post
// query contains id of post
func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	idParam := r.URL.Query().Get("id")

	log.Println(idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err := api.db.PostDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
