package structures

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

// =============================================
// WS TRANSFER OBJECTS
// =============================================
type SongPayload struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	Source        string `json:"source"`
	Art           string `json:"art"`
	Album         string `json:"album"`
	Id            int    `json:"id"`
	TotalLengthMs int64
}

type ClientPayload struct {
	Id          int    `json:"id"`
	ChannelId   int    `json:"channel_id"`
	DisplayName string `json:"display_name"`
}

type ChannelMessagePayload struct {
	Message  string `json:"message"`
	SenderId int    `json:"sender_id"`
}

type ChannelPayload struct {
	JoinTimestamp int64       `json:"join_timestamp"`
	JoinClientId  int         `json:"join_client_id"`
	Name          string      `json:"name"`
	NowPlaying    SongPayload `json:"now_playing"`
	NumUsers      int         `json:"num_users"`
	Id            int         `json:"id"`
}

type TimeSyncPayload struct {
	ChannelId    int   `json:"channel_id"`
	SongLocation int64 `json:"song_location"`
	ServerTime   int64 `json:"server_time"`
}

type ErrorMessage struct {
	CausedByType string `json:"caused_by_type"`
	Message      string `json:"message"`
}

// =============================================
// SERVER STRUCTURES
// =============================================

type ConnectionObject struct {
	Socket websocket.Conn
	Client ClientPayload
}

type ChannelObject struct {
	Channel       ChannelPayload
	SongStartTime int64
	Song          SongPayload
	Users         []ConnectionObject
}

type SinglePageAppHandler struct {
	StaticPath string
	IndexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h SinglePageAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	path = filepath.Join(h.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	log.Printf("Checking for: %s", path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
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
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}

// =============================================
// UTILITY METHODS
// =============================================

//Returns the current time since epoch in ms [Javascript Date() format]
func NowInMs() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond))
}
