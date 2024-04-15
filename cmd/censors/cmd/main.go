package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"sf-mu/cmd/censors/pkg/api"
	"sf-mu/cmd/censors/pkg/supply"

	config "sf-mu/pkg/configs"
	"sf-mu/pkg/middleware"
	"sf-mu/pkg/models/censors"
	censorsstorage "sf-mu/pkg/storage/censors"
	"time"

	"github.com/joho/godotenv"
)

type server struct {
	db  censors.Interface
	api *api.API
}

func init() {
	// loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	var srv server

	cfg := config.New()

	dbURL := cfg.Comments.URLdb
	//default port
	port := cfg.Censor.AdrPort

	portFlag := flag.String("censor-port", port, "censor api port")
	imFlag := flag.Bool("i", false, "install mode")
	flag.Parse()

	im := false
	if imFlag != nil {
		im = *imFlag
	}

	portCensor := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//
	db, err := censorsstorage.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if im {
		log.Println("install mode activated")

		log.Println("drop black list table")
		// Drop blacklist if exists
		err = db.DropBlackListTable()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("create black list table")
		// create table blacklist
		err = db.CreateBlackListTable()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("load black list from file")
		//load words.txt into blacklist table
		blackList, err := supply.BlackList()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(blackList)
			log.Println("insert loaded black list into black list table")
			for _, v := range blackList {
				err := db.AddList(v)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	// init db
	srv.db = db
	srv.api = api.New(srv.db)
	srv.api.Router().Use(middleware.Middle)
	log.Print("listen http://localhost" + portCensor)

	err = http.ListenAndServe(portCensor, srv.api.Router())
	if err != nil {
		log.Fatal("Failed to start server. Error:", err)
	}

}
