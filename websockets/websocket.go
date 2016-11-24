package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"
	"log"
)

func WSEcho(ws *websocket.Conn) {
	for {
		var reply string

		err := websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't Receive : ",err.Error())
			break
		}

		fmt.Println("Recevied from client : ",reply)

		msg := "Server : "+reply

		err = websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't Send : ",err.Error())
			break
		}
	}
}

func main() {
	http.Handle("/",websocket.Handler(WSEcho))

	log.Fatal(http.ListenAndServe(":1234", nil))
}