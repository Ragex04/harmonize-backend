package main

import (
	"encoding/json"
	"fmt"
	"harmonize-server/structures"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const PORT int = 8000

var upgrader = websocket.Upgrader{}

func createNewWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("New WS Requested!!")
	// allow anyone lol
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	ws.Close()
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
