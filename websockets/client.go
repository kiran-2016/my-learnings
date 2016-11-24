package main

import (
	"os"
	"golang.org/x/net/websocket"
	"bufio"
	"log"
	"time"
	"flag"
	"strconv"
)

func send_to_server(conn *websocket.Conn, stream chan interface{}) {
	for {
		select {
		case msg := <-stream:
			err := websocket.Message.Send(conn, msg)
			if err != nil {
				log.Println("Error sending message to server : ", err.Error())
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func receive_from_server(conn *websocket.Conn) {
	var msg string
	remote_address := conn.RemoteAddr().String()

	for {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Println("Error Receiving message from server : ", err.Error())
			return
		}

		log.Printf("Received message from server %s: %s", remote_address, msg)
	}
}

func clean_up(quit chan int, exit_id int) {
	for {
		if _, err := os.Stat("/tmp/cexit"+strconv.Itoa(exit_id)); err == nil {
			quit <- exit_id
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func long_text(service string, quit chan int, exit_id int, forever bool) {
	fname := "/tmp/dict.txt"

	conn, err := websocket.Dial(service, "", "http://localhost")
	if err != nil {
		log.Println("Error connecting to server : ", err.Error())
		return
	}

	defer conn.Close()

	go clean_up(quit, exit_id)

	for forever {
		fin, err := os.Open(fname)
		if err != nil {
			log.Println("Unable to open file : ", err)
			return
		}

		stream := make(chan interface{})

		go send_to_server(conn, stream)
		go receive_from_server(conn)

		reader := bufio.NewScanner(fin)
		for reader.Scan() {
			stream <- reader.Text()
		}

		fin.Close()
	}
}

func two_way_comm_client(service string, quit chan int, exit_id int, forever bool) {
	conn, err := websocket.Dial(service, "", "http://localhost")
	if err != nil {
		log.Println("Error connecting to server : ", err.Error())
		return
	}

	defer conn.Close()

	stream := make(chan interface{})

	go send_to_server(conn, stream)
	go receive_from_server(conn)

	go clean_up(quit, exit_id)

	for forever {
		stream <- "hello from client"
		time.Sleep(2 * time.Second)
	}
}

func main() {
	url := flag.String("url", "ws://localhost:8080/long_hello", "Specify websocket URL to connect to ws://<host>:<port>/path")
	thread_count := flag.Int("threads", 1, "Number of web socket clients to initiate!")
	forever := flag.Bool("forever", true, "By default it runs for ever!")

	flag.Parse()

	quit := make(chan int)
	for i := 1; i <= *thread_count; i++ {
		//go long_text(*url, quit, i, *forever)
		go two_way_comm_client(*url, quit, i, *forever)
	}

	for i := 1; i <= *thread_count; i++ {
		<- quit
	}

}