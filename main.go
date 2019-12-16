package main

import (
	"encoding/json"
	"fmt"
	"harmonize-server/connections"
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
		make(chan map[string]interface{}, 10),
		false,
		false,
	}
	log.Printf("New WS Requested!! Added client: %d", con.Client.Id)
	go connections.HandleNewConnection(&con)
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
