package connections

import (
	"encoding/json"
	"harmonize-server/structures"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const HOST string = "http://localhost:8000/"
const MAX_USERS_PER_CHANNEL = 10

var CHANNELS *[]structures.ChannelObject

func HandleSafeSend(con *structures.ConnectionObject) {
	for data := range con.Txd {
		con.Socket.WriteJSON(&data)
	}
	log.Printf("%d: Sending thread ending...", con.Client.Id)
}

func ReadChannels() {
	f, err := os.Open("./data/channels.json")
	if err != nil {
		log.Printf("Error opening channels data: %v", err)
	}

	byteData, _ := ioutil.ReadAll(f)

	defer f.Close()

	var channels []structures.ChannelObject

	json.Unmarshal(byteData, &channels)

	log.Printf("%v", channels)

	CHANNELS = &channels
}

func GetChannelByID(id int) *structures.ChannelObject {

	for _, channel := range *CHANNELS {
		if channel.Channel.Id == id {
			return &channel
		}
	}
	return nil
}

func HandleNewConnection(con *structures.ConnectionObject) {
	ReadChannels()
	go HandleSafeSend(con)
	for {
		_, msg, err := con.Socket.ReadMessage()
		if err != nil {
			log.Printf("%d: Error reading from ws from client", con.Client.Id)
			con.Disconnected = true
			con.Socket.Close()
			close(con.Txd)
			break
		}

		var msgType structures.BasicMessage
		json.Unmarshal(msg, &msgType)

		var typeaction []string
		var mtype, action string
		typeaction = strings.Split(string(msgType.MsgType), "/")

		mtype = typeaction[0]
		action = typeaction[1]

		log.Printf("%d: Got %s:%s", con.Client.Id, mtype, action)

		switch mtype {
		case "CHANNEL":
			//handleChannelMessage(con, msg, action)
			go handleChannelMessage(con, msg, action)
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
	cnn.Id = c.Id
	cnn.Offset = structures.NowInMs() - c.JoinTimestamp
	cnn.VoteOptions = make([]structures.SongPayload, 1)

	var theChannel *structures.ChannelObject
	theChannel = GetChannelByID(c.Id)

	if theChannel == nil {
		log.Printf("Cannot find channel with ID: %d", c.Id)
		return
	}

	yeet := make(map[string]interface{})

	yeet["type"] = "CHANNEL/INFO"
	yeet["payload"] = cnn

	con.Txd <- yeet

	if theChannel.SongStartTime == 0 {
		theChannel.SongStartTime = structures.NowInMs()
		theChannel.Users = make([]structures.ConnectionObject, MAX_USERS_PER_CHANNEL)
	}

	//theChannel.Users = append(theChannel.Users, *con)
	userDuplicate := false

	for _, user := range theChannel.Users {
		if user.Client.Id == con.Client.Id {
			userDuplicate = true
		}
	}

	if !userDuplicate {
		theChannel.Users = append(theChannel.Users, *con)
		go StartSyncPackets(c.Id, con)
	} else {
		log.Printf("Duplicate User Found!!")
	}

}

func StartSyncPackets(channelId int, con *structures.ConnectionObject) {
	channel := GetChannelByID(channelId)

	for {
		time.Sleep(100 * time.Millisecond)
		pkt := &structures.TimeSyncPayload{
			channelId,
			(structures.NowInMs() - channel.SongStartTime) % channel.Song.TotalLengthMs,
			structures.NowInMs(),
		}

		data := make(map[string]interface{})
		data["type"] = "SYNC/PACKET"
		data["payload"] = pkt

		//bin_data, _ := json.Marshal(data)

		con.Txd <- data
	}
}
