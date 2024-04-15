package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"sf-mu/cmd/news/pkg/api"
	"sf-mu/cmd/news/pkg/rss"
	config "sf-mu/pkg/configs"
	"sf-mu/pkg/middleware"
	"sf-mu/pkg/models/news"
	newsstorage "sf-mu/pkg/storage/news"

	"time"

	"github.com/joho/godotenv"
)

type server struct {
	db  news.Interface
	api *api.API
}

func init() {
	// loads environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	configURL = "./cmd/news/cmd/config.json"
)

func main() {

	var srv server

	cfg := config.New()
	dbURL := cfg.News.URLdb
	// default port
	port := cfg.News.AdrPort
	portFlag := flag.String("news-port", port, "news server port")
	imFlag := flag.Bool("i", false, "install mode")
	flag.Parse()
	im := false
	if imFlag != nil {
		im = *imFlag
	}
	portNews := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := newsstorage.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if im {
		log.Println("install mode activated")
		log.Println("drop news table")
		err = db.DropGonewsTable()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("create news table")
		err = db.CreateGonewsTable()
		if err != nil {
			log.Println(err)
			return
		}

	}
	srv.db = db

	srv.api = api.New(srv.db)

	//--------------------------------------------------------
	log.Println("making channels")
	chanPosts := make(chan []news.Post)
	chanErrs := make(chan error)

	go func() {
		err := rss.GoNews(configURL, chanPosts, chanErrs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for posts := range chanPosts {
			if err := srv.db.PostsCreation(posts); err != nil {
				chanErrs <- err
			}
		}
	}()

	go func() {
		for err := range chanErrs {
			log.Println(err)
		}
	}()

	srv.api.Router().Use(middleware.Middle)

	log.Print("Starting server on http://localhost" + portNews)

	err = http.ListenAndServe(portNews, srv.api.Router())
	if err != nil {
		log.Fatal("Failed to start server. Error:", err)
	}

}
