package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
	"github.com/rs/cors"
	"github.com/Zhanat87/go/apis"
)

var (
	Env string
)
// cd ~/go/socketio-server && go build -ldflags "-X main.Env=docker" -o ~/go/bin/socketio-server
// @link http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
func main() {
	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to room
		c.Join("golang")
	})

	//handle custom event
	server.On("socialAuth", func(c *gosocketio.Channel, msg apis.SocialAuthMessage) string {
		c.BroadcastTo("golang", "socialAuth" + msg.Uuid, msg)
		return "OK"
	})

	//setup http server
	mux := http.NewServeMux()

	var frontendUrl string
	var backendUrl string
	if Env == "docker" {
		frontendUrl = "http://zhanat.site:8081"
		backendUrl = "http://zhanat.site:8080"
	} else {
		frontendUrl = "http://localhost:3000"
		backendUrl = "http://localhost:8080"
	}
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{frontendUrl, backendUrl},
		AllowCredentials: true,
	}).Handler(mux)

	mux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":5000", handler))
}
