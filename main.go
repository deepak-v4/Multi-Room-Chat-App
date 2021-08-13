package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan Msg)

// Configure upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Define message

type Msg struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	RoomID   string `json:"roomid"`
}

func main() {

	//file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	//config websocket route
	http.HandleFunc("/ws", handleConnections)

	//handle incoming msg
	go handleMsg()

	//start server

	log.Println("Starting server on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Listen and serve:", err)
	}

}

func handleConnections(wr http.ResponseWriter, rd *http.Request) {

	//upgrade initial request to a websocket
	w, err := upgrader.Upgrade(wr, rd, nil)
	if err != nil {
		log.Fatal(err)
	}

	//close the connection when function return
	defer w.Close()

	var initmsg Msg
	initerr := w.ReadJSON(&initmsg)

	if initerr != nil {
		log.Printf("error:%v", initerr)
		return
	}

	//register new client
	room_id := initmsg.RoomID
	inner, ok := clients[room_id]
	if !ok {
		inner = make(map[*websocket.Conn]bool)
		clients[room_id] = inner
	}
	inner[w] = true

	//fmt.Println("assignment to map done")
	broadcast <- initmsg

	for {
		var message Msg
		//fmt.Println("before msg read")
		//read new message as JSON and map it to Msg obj
		err := w.ReadJSON(&message)
		if err != nil {
			log.Printf("error:%v", err)
			delete(clients[room_id], w)
			break
		}

		// send newly received msg to broadcast channel
		broadcast <- message

	}

}

func handleMsg() {

	for {

		//fmt.Println("waiting for broadcast")
		//grab the next msg from the broadcast channel
		message := <-broadcast

		room_id := message.RoomID

		for client := range clients[room_id] {
			//fmt.Println("user msg ::", message)

			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error:%v", err)
				client.Close()
				delete(clients[room_id], client)
			}
		}

	}

}
