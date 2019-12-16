package main

import (
	"encoding/json"
	"fmt"
	"harmonize-server/structures"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const PORT int = 8000

var upgrader = websocket.Upgrader{}

func RecvMsgs(con *structures.ConnectionObject) {
	for {
		_, msg, err := con.Socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading from ws from client: %d", con.Client.Id)
		}
		log.Printf("Got: %s", string(msg))
		con.Recvd <- msg

	}
}

func handleConnection(con *structures.ConnectionObject) {
	go RecvMsgs(con)
}

func createNewWebsocket(w http.ResponseWriter, r *http.Request) {

	// allow anyone lol
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	con := structures.ConnectionObject{
		ws,
		structures.ClientPayload{
			rand.Intn(int(1000000)),
			-1,
			"TmpClient",
		},
		make(chan []byte, 10),
		false,
	}
	log.Printf("New WS Requested!! Added client: %d", con.Client.Id)
	handleConnection(&con)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/ws", createNewWebsocket)

	spa := structures.SinglePageAppHandler{
		StaticPath: "./wwwroot",
		IndexPath:  "index.html",
	}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", PORT),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on port: %d", PORT)
	log.Fatal(srv.ListenAndServe())
}
