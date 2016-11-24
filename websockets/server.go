package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"flag"
	"strconv"
	"time"
)

func hey(conn *websocket.Conn) {
	defer conn.Close()

	remote_address := conn.RemoteAddr().String()
	err := websocket.Message.Send(conn, "Hey, "+remote_address+"\n")
	if err != nil {
		log.Println("Error sending message : ", err.Error())
	}

	time.Sleep(1 * time.Second)
}

func hello(conn *websocket.Conn) {
	var reply string
	address := conn.LocalAddr().String()

	for {
		err := websocket.Message.Receive(conn, &reply)
		if err != nil {
			log.Println("Error receiving message : ", err.Error())
			break
		}

		err = websocket.Message.Send(conn, "from " + address + " : " + reply)
		if err != nil {
			log.Println("Error sending message : ", err.Error())
			break
		}
	}
}

func long_hello(conn *websocket.Conn) {
	defer conn.Close()

	address := conn.LocalAddr().String()
	remote_address := conn.Request().RemoteAddr

	log.Printf("Connected to client : %s", remote_address)

	for {
		var msg string

		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Printf("Error Receiving message from client (%s): %s", remote_address, err.Error())
			break
		}

		pmsg := fmt.Sprintf("from %s to %s : ",address, remote_address)
		err = websocket.Message.Send(conn, pmsg+msg)
		if err != nil {
			log.Printf("Error sending message to client (%s) : %s", remote_address, err.Error())
			break
		}
	}
}

func two_way_comm_server(conn *websocket.Conn) {
	defer conn.Close()

	remote_address := conn.Request().RemoteAddr

	quit := make(chan int)

	go func() {
		for {
			var msg string

			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				log.Printf("Error Receiving message from client (%s): %s", remote_address, err.Error())
				break
			}

			log.Printf("Received message from %s : %s", remote_address, msg)
		}
	}()

	go func() {
		msg := "hello from server"
		for {
			err := websocket.Message.Send(conn, msg)
			if err != nil {
				log.Printf("Error sending message to client (%s) : %s", remote_address, err.Error())
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	<- quit
}

func init_server(port int) {
	address := ":"+strconv.Itoa(port)

	mux := http.NewServeMux()
	server := &http.Server{Addr:address, Handler:mux}

	mux.Handle("/hello", websocket.Handler(hello))
	mux.Handle("/long_hello", websocket.Handler(long_hello))
	mux.Handle("/two_way", websocket.Handler(two_way_comm_server))
	mux.Handle("/hey", websocket.Handler(hey))
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Welcome to websocket programming!\n"))
		rw.Write([]byte(fmt.Sprintf("Server running at port : %d\n", port)))
	})

	log.Println("Starting server at port : ", port)

	server.ListenAndServe()
}

func init_secure_server(port int) {
	address := ":"+strconv.Itoa(port)

	mux := http.NewServeMux()
	server := &http.Server{Addr:address, Handler:mux}

	mux.Handle("/hello", websocket.Handler(hello))
	mux.Handle("/long_hello", websocket.Handler(long_hello))
	mux.Handle("/two_way", websocket.Handler(two_way_comm_server))
	mux.Handle("/hey", websocket.Handler(hey))
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Welcome to websocket (secure) programming!\n"))
		rw.Write([]byte(fmt.Sprintf("Server running at port : %d\n", port)))
	})

	log.Println("Starting server(secure) at port : ", port)

	log.Fatal(server.ListenAndServeTLS("server.pem", "server.key"))
}

func main() {
	port := flag.Int("port", 8080, "Specify port to run server!")
	thread_count := flag.Int("threads", 1, "Specify number of servers to spawn")
	secure_port := flag.Int("secure-port", 8443, "Specify secure port to run server!")

	flag.Parse()
	quit := make(chan int)

	for i := 0; i < *thread_count; i++ {
		go init_server(*port+i)
		go init_secure_server(*secure_port+i)
	}

	<- quit
}
