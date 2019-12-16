package structures

import "time"

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

`Returns the current time since epoch in ms [Javascript Date() format]`
func NowInMs() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond))
}
