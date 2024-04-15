package main

import (
	"flag"
	"sf-mu/cmd/gateway/pkg/api"

	config "sf-mu/pkg/configs"
	"sf-mu/pkg/middleware"

	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// gateway server
type server struct {
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

	//default ports
	port := cfg.Gateway.AdrPort
	newsPort := cfg.News.AdrPort
	censorPort := cfg.Censor.AdrPort
	comment := cfg.Comments.AdrPort

	portFlag := flag.String("gateway-port", port, "gateway port")
	portFlagNews := flag.String("news-port", newsPort, "news server port")
	portFlagCensor := flag.String("censor-port", censorPort, "censor server port")
	portFlagComment := flag.String("comments-port", comment, "comments server port")

	flag.Parse()

	portGateway := *portFlag
	portNews := *portFlagNews
	portCensor := *portFlagCensor
	portComment := *portFlagComment

	srv.api = api.New(cfg, portNews, portCensor, portComment)
	srv.api.Router().Use(middleware.Middle)

	log.Print("Starting server on http://localhost" + portGateway + "/news")

	err := http.ListenAndServe(portGateway, srv.api.Router())
	if err != nil {
		log.Fatal("Failed to start server. Error:", err)
	}

}
