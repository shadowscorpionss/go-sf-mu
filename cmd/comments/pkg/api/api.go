package api

import (
	"encoding/json"
	"net/http"
	comments "sf-mu/pkg/models/comments"
	"strconv"

	"github.com/gorilla/mux"
)


type API struct {
	r  *mux.Router        // api Router
	db comments.Interface // 
}

// creates comments api
func New(db comments.Interface) *API {
	api := API{
		r:  mux.NewRouter(),
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// returns router
func (api *API) Router() *mux.Router {
	return api.r
}

// creates api endpoints
func (api *API) endpoints() {


	api.r.HandleFunc("/comments", api.commentsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments/del", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)
}

// commentsHandler, responses with news_id comments 
// query contains news_id
func (api *API) commentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	parseId := r.URL.Query().Get("news_id")

	newsId, err := strconv.Atoi(parseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comments, err := api.db.AllComments(newsId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// adds comment
func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c comments.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = api.db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.ResponseWriter.WriteHeader(w, http.StatusCreated)
}

// deletes comment
func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var c comments.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
