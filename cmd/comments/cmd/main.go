package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"sf-mu/cmd/comments/pkg/api"
	config "sf-mu/pkg/configs"
	"sf-mu/pkg/middleware"
	"sf-mu/pkg/models/comments"
	commentsstorage "sf-mu/pkg/storage/comments"
	"time"

	"github.com/joho/godotenv"
)

type server struct {
	db  comments.Interface
	api *api.API
}

func init() {
	// loads environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	var srv server

	cfg := config.New()

	dbURL := cfg.Comments.URLdb
	// default port
	port := cfg.Comments.AdrPort

	portFlag := flag.String("comments-port", port, "comments server port")
	imFlag := flag.Bool("i", false, "install mode")
	flag.Parse()
	im := false
	if imFlag != nil {
		im = *imFlag
	}

	portComments := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := commentsstorage.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if im {
		log.Println("install mode activated")
		log.Println("drop comment table")
		err = db.DropCommentTable()
		if err != nil {
			log.Println(err)
			return
		}
		
		log.Println("create comment table")
		err = db.CreateCommentTable()
		if err != nil {
			log.Println(err)
			return
		}
	}

	srv.db = db

	srv.api = api.New(srv.db)

	srv.api.Router().Use(middleware.Middle)

	log.Print("Starting server on http://localhost" + portComments)

	err = http.ListenAndServe(portComments, srv.api.Router())
	if err != nil {
		log.Fatal("Failed to start server. Error:", err)
	}
}
