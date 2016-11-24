package main

import (
	"golang.org/x/net/websocket"
	"github.com/prometheus/common/log"
	"io"
	"fmt"
	"flag"
)

func insecure_call(url string) {
	conn, err := websocket.Dial(url, "", "http://localhost:8080")
	if err != nil {
		log.Fatal("Unable to connect : ", err.Error())
	}

	defer conn.Close()

	var msg string = "Hello from client!"
	count := 1
	for count > 0 {
		err = websocket.Message.Send(conn, msg)
		if err != nil {
			log.Fatal("Error sending message : ", err.Error())
		}

		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal("Error receiving message : ", err.Error())
		}

		fmt.Println("Message received from server : ", msg)
		count--
	}
}

func secure_call(url string) {

}

func main() {
	secure := flag.Bool("secure", false, "Specify whether http/https connection!")
	url := flag.String("url", "ws://localhost:8080/hello", "Specify URL of websoket server!")

	if secure {
		url = "wss://localhost:8080/hello"
	}

	quit := make(chan int)

	if secure {
		go secure_call(url)
	} else {
		go insecure_call(url)
	}

	<- quit
}
