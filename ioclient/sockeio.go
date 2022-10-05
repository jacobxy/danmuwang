package ioclient

import (
	"log"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var Server *gosocketio.Server

const (
	ROOM_NAME = "home"
)

func BroadCast(key string, b interface{}) error {
	GetSocketIO().BroadcastTo(ROOM_NAME, key, b)
	return nil
	channel, err := GetSocketIO().GetChannel(ROOM_NAME)
	if err != nil {
		return err
	}
	return channel.Emit(key, b)
}

func GetSocketIO() *gosocketio.Server {
	if Server != nil {
		return Server
	}
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		c.Join(ROOM_NAME)
		c.BroadcastTo("chat", "message", "xxx")

		//of course, you can list the clients in the room, or account them
		//or check the amount of clients in room
		log.Println(c.Id())
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		//caller is not necessary, client will be removed from rooms
		//automatically on disconnect
		//but you can remove client from room whenever you need to

		log.Println("Disconnected")
	})
	server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("Error occurs")
	})

	type Message struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	//handle custom event
	server.On("send", func(c *gosocketio.Channel, msg Message) string {
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	server.On("say", func(c *gosocketio.Channel, msg Message) string {
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	Server = server
	return server
	//setup http server
}
