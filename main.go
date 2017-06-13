package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
	"github.com/rs/cors"
	"github.com/Zhanat87/go/apis"
	"fmt"
)

var (
	Env string
)

type ChatMessage struct {
	ChatUser
	Message  string `json:"message"`
	Time     string `json:"time"`
}

type ChatUser struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func main() {
	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		// join them to room
		c.Join("socialAuthRoom")

		// join them to chat room
		c.Join("chatRoom")

		//for a, b := range c.RequestHeader() {
		//	log.Println("key: " + a)
		//	for c, d := range b {
		//		log.Println(fmt.Sprintf("%d: %s", c, d))
		//	}
		//}
	})

	// on disconnection handler, if client hangs connection unexpectedly, it will still occurs
	// you can omit function args if you do not need them
	// you can return string value for ack, or return nothing for emit
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		// caller is not necessary, client will be removed from rooms
		// automatically on disconnect
		// but you can remove client from room whenever you need to
		c.Leave("socialAuthRoom")
		c.Leave("chatRoom")

		log.Println("Disconnected")
	})

	// error catching handler
	server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("Error occurs")
	})

	// handle custom event
	server.On("socialAuth", func(c *gosocketio.Channel, msg apis.SocialAuthMessage) string {
		c.BroadcastTo("socialAuthRoom", "socialAuth" + msg.Uuid, msg)
		return "OK"
	})

	// handle custom event
	server.On("chatMessage", func(c *gosocketio.Channel, msg ChatMessage) string {
		// send event to all in room
		c.BroadcastTo("chatRoom", "chatMessage", msg)
		return "OK from chatMessage"
	})

	// handle custom event
	server.On("chatUsers", func(c *gosocketio.Channel, msg ChatUser) string {
		/*
		send event to all in room
		@link http://stackoverflow.com/questions/10058226/send-response-to-all-clients-except-sender-socket-io
		note: sending to all clients in 'chatUsers' room(chatRoom) except sender, not work on golang backend
		 */
		c.BroadcastTo("chatRoom", "chatUsers", msg)
		log.Println(msg)
		return "OK from chatUsers"
	})

	// handle custom event
	server.On("forceDisconnect", func(c *gosocketio.Channel, msg string) string {
		c.Leave("socialAuthRoom")
		c.Leave("chatRoom")

		log.Println("forceDisconnect")

		return "OK from forceDisconnect"
	})

	// handle custom event
	server.On("chatUserLogout", func(c *gosocketio.Channel, userId int) string {
		c.Leave("socialAuthRoom")
		c.Leave("chatRoom")

		log.Println(fmt.Sprintf("chatUserLogout: %d", userId))
		c.BroadcastTo("chatRoom", "chatUserLogout", userId)

		return "OK"
	})

	// setup http server
	mux := http.NewServeMux()

	//var frontendUrl string
	//var backendUrl string
	//if Env == "docker" {
	//	frontendUrl = "http://zhanat.site:8081"
	//	backendUrl = "http://zhanat.site:8080"
	//} else {
	//	frontendUrl = "http://localhost:3000"
	//	backendUrl = "http://localhost:8080"
	//}
	allowedOrigins := []string{"http://zhanat.site:8081", "http://zhanat.site:8080",
		"http://localhost:3000", "http://localhost:8080"}
	handler := cors.New(cors.Options{
		//AllowedOrigins: []string{frontendUrl, backendUrl},
		AllowedOrigins: allowedOrigins,
		AllowCredentials: true,
	}).Handler(mux)

	mux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":5000", handler))
}
