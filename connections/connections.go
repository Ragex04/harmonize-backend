package connections

import (
	"encoding/json"
	"harmonize-server/structures"
	"log"
	"strings"
)

const HOST string = "http://localhost:8000/"

func HandleNewConnection(con *structures.ConnectionObject) {
	for {
		_, msg, err := con.Socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading from ws from client: %d", con.Client.Id)
		}

		var msgType structures.BasicMessage
		json.Unmarshal(msg, &msgType)

		var typeaction []string
		var mtype, action string
		typeaction = strings.Split(string(msgType.MsgType), "/")

		mtype = typeaction[0]
		action = typeaction[1]

		log.Printf("Got %s:%s", mtype, action)

		switch mtype {
		case "CHANNEL":
			handleChannelMessage(con, msg, action)
			//go handleChannelMessage(con, msg, action)
		default:
			log.Printf("ERROR: No messgae handler of type: %s", mtype)
		}
	}
}

func handleChannelMessage(con *structures.ConnectionObject, data []byte, action string) {
	var ChannelMessage structures.ChannelPayload
	json.Unmarshal(data, &ChannelMessage)

	switch action {
	case "JOIN":
		handleChannelJoin(con, ChannelMessage)
	default:
		log.Printf("ERROR: No CHANNEL handler of action type: %s", action)
	}
}

func handleChannelJoin(con *structures.ConnectionObject, c structures.ChannelPayload) {

	s := structures.SongPayload{
		"Movement",
		"Oliver Tree",
		HOST + "audio/movement.mp3",
		"https://i1.sndcdn.com/artworks-000347035692-n5238t-t500x500.jpg",
		"Movement",
		1,
		163000,
	}

	var cnn structures.ChannelPayload
	cnn.NowPlaying = s
	cnn.Id = 1
	cnn.Offset = structures.NowInMs() - c.JoinTimestamp
	log.Printf("Sending join msg")

	yeet := make(map[string]interface{})

	//payload, _ := json.Marshal(cnn)
	yeet["type"] = "CHANNEL/INFO"
	yeet["payload"] = cnn

	con.Socket.WriteJSON(&yeet)

}
