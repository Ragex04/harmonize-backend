package main

import (
	"encoding/json"
	"log"
	"net/http"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.

var upgrader = websocket.Upgrader{}

type spaHandler struct {
	staticPath string
	indexPath  string
}

type SourcesThing struct {
	Enabled bool `json:"enabled"`
	Url string `json:"url"`
	Pan int `json:"pan"`
	Gain int `json:"gain"`
}

type TimeMsg struct {
	MsgType    string `json:"type"`
	Timecode int64 `json:"timecode"`
	Timestamp  int64 `json:"timestamp"`
	Offset int64 `json:"offset"`
	Name int `json:"name"`
	Id int `json:"id"`
	Sources []SourcesThing `json:"sources"`
	Url string `json:"url"`
	Enabled bool `json:"enabled"`
}

var startTime int64 = -1

var clients = make(chan *websocket.Conn, 10)

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
        // if we failed to get the absolute path respond with a 400 bad request
		// and stop
		log.Printf("Some sort of bad abs path: %s", r.URL.Path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

    // check whether a file exists at the given path
	_, err = os.Stat(path)
	log.Printf("Checking for: %s", path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
        // if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		log.Printf("File not found: %s", path)
		log.Printf("Err: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	log.Printf("\tServing: %s", path)
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func handleWSConnection(){
	audioLength := 200782
	var clientID int
	ws := <- clients
	var msg TimeMsg
	// Read in a new message as JSON and map it to a Message object
	err := ws.ReadJSON(&msg)
	if err != nil {
			log.Printf("error: %v", err)
			//delete(clients, ws)
	}
	// Send the newly received message to the broadcast channel
	log.Printf("\tGot new WS Message!!: %v",msg)
	
	if msg.MsgType == "hello"{
		if startTime == -1 {
			startTime = (time.Now().UnixNano() / int64(time.Millisecond))
		}
		clientID = rand.Intn(int(1000000))
		var outmsg1, outmsg2 TimeMsg
		outmsg1.MsgType = "hello"
		outmsg1.Id = clientID
		outmsg1.Name = clientID
		outmsg1.Offset = (time.Now().UnixNano() / int64(time.Millisecond)) - msg.Timestamp
		outmsg1.Timestamp = (time.Now().UnixNano() / int64(time.Millisecond))
		ws.WriteJSON(outmsg1)

		outmsg2.MsgType = "set_sources"
		outmsg2.Id = clientID
		outmsg2.Sources = []SourcesThing{SourcesThing{true, "audio/test.ogg", 0, 1}}
		ws.WriteJSON(outmsg2)

		outmsg2.Name = clientID
		ws.WriteJSON(outmsg2)


		go func(){
			log.Printf("Starting client in 5 seconds....")
			time.Sleep(5*time.Second)
			log.Printf("Client %d started", clientID)
			var startMsg, setMsg TimeMsg
			startMsg.Id = clientID
			startMsg.MsgType = "admin_sources"
			startMsg.Enabled = true
			startMsg.Url = "audio/test.ogg"
			ws.WriteJSON(startMsg)

			setMsg.MsgType = "set_sources"
			setMsg.Id = clientID
			setMsg.Sources = []SourcesThing{SourcesThing{true, "audio/test.ogg", 0, 1}}
			ws.WriteJSON(setMsg)
		}()

		for {
			var msg3 TimeMsg
			msg3.Id = clientID
			msg3.Timestamp = (time.Now().UnixNano() / int64(time.Millisecond))
			msg3.MsgType = "timecode"
			msg3.Timecode = ((time.Now().UnixNano() / int64(time.Millisecond)) - startTime) % int64(audioLength)
			ws.WriteJSON(msg3)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// msg.MsgType = "TEST"
	// msg.Timecode = -1
	// msg.Timestamp = -2
	// ws.WriteJSON(msg)
	defer ws.Close()
}

func getWS(w http.ResponseWriter, r *http.Request){
	log.Println("New WS Connection!")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Fatal(err)
		}
		
	clients <- ws
	go handleWSConnection()
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/ws", getWS)

	spa := spaHandler{staticPath: "../web-audio-stream-sync/frontend/build/", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}