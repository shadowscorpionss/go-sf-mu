package api

import (
	"bytes"
	"encoding/json"
	config "sf-mu/pkg/configs"

	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

// Application PI
type API struct {
	r           *mux.Router // application router
	cfg         *config.Config
	portNews    string
	portCensor  string
	portComment string
}

// Constructor of api
func New(cfg *config.Config, portNews, portCensor, portComment string) *API {
	api := API{
		r:           mux.NewRouter(),
		cfg:         cfg,
		portNews:    portNews,
		portCensor:  portCensor,
		portComment: portComment,
	}
	api.endpoints()
	return &api
}

type ResponseDetailed struct {
	NewsDetailed struct {
		ID      int    `json:"ID"`
		Title   string `json:"Title"`
		Content string `json:"Content"`
		PubTime int    `json:"PubTime"`
		Link    string `json:"Link"`
	} `json:"NewsDetailed"`
	Comments []struct {
		ID              int    `json:"ID"`
		NewsID          int    `json:"newsID"`
		ParentCommentID int    `json:"parentCommentID"`
		Content         string `json:"content"`
		PubTime         int    `json:"pubTime"`
	} `json:"Comments"`
}

// returns api router
func (api *API) Router() *mux.Router {
	return api.r
}

// sets up api endpoints
func (api *API) endpoints() {
	api.r.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/latest", api.newsLatestHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/search", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments/del", api.delCommentHandler).Methods(http.MethodDelete, http.MethodOptions)

}

// responses with list of posts
func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news" {
		http.NotFound(w, r)
	}
	portNews := api.portNews

	//
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost" + portNews, //
	})

	//
	proxy.ServeHTTP(w, r)

}

func (api *API) newsLatestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news/latest" {
		http.NotFound(w, r)
	}
	portNews := api.portNews

	//
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost" + portNews, //
	})

	//
	proxy.ServeHTTP(w, r)

}

func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news/search" {
		http.NotFound(w, r)
	}
	portNews := api.portNews
	portComment := api.portComment

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "query params are to be set", http.StatusBadRequest)
		return
	}

	chNews := make(chan *http.Response, 2)
	chComments := make(chan *http.Response, 2)
	chErr := make(chan error, 2)
	var response ResponseDetailed
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		//
		resp1, err := http.Get("http://localhost" + portNews + "/news/search" + "?id=" + idParam)
		chErr <- err
		chNews <- resp1
	}()

	go func() {
		defer wg.Done()
		//
		resp2, err := http.Get("http://localhost" + portComment + "/comments" + "?news_id=" + idParam)
		chErr <- err
		chComments <- resp2
	}()

	wg.Wait()
	close(chErr)

	for err := range chErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

block:
	for {
		select {
		case respNews := <-chNews:
			body, err := io.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := io.ReadAll(respComment.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.Comments)
		default:
			break block
		}

	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/add" {
		http.NotFound(w, r)
	}
	portCensor := api.portCensor
	portComment := api.portComment

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	Body1 := io.NopCloser(bytes.NewBuffer(bodyBytes))
	Body := io.NopCloser(bytes.NewBuffer(bodyBytes))

	respCensor, err := MakeRequest(r, http.MethodPost, "http://localhost"+portCensor+"/comments/check", Body1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "wrong comment content", respCensor.StatusCode)
		return
	}

	resp, err := MakeRequest(r, http.MethodPost, "http://localhost"+portComment+"/comments/add", Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (api *API) delCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/del" {
		http.NotFound(w, r)
	}
	portComment := api.portComment

	//
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost" + portComment, //
	})

	//
	proxy.ServeHTTP(w, r)
}

func MakeRequest(req *http.Request, method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Header
	return client.Do(r)
}
