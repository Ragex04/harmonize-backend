class Song {
    type: string; // "SONG"
    action: string // "VOTE" || "GETALL"
    title: string; // "White Ash"
    artist: string; // "The Pillows"
    source: string; // "http://host/song.mp3"
    art: string; // "http://hpst/art/FCFL.jpg"
    album: string // "FCFL"
    id: number; // 253
}

class Client {
    type: string; // "CLIENT"
    action: string; // "JOIN" || "GETINFO"
    id: number; // 69
    channel_id: number; // 12
    display_name: string; // "Elon Moosk"
}

class ChannelMessage {
    type: string; // "MESSAGE"
    sender_id: number;
    message: string;
}

class Channel {
    type: string; // "CHANNEL"
    action: string // "JOIN" || "GETINFO" || "LEAVE"
    join_timestamp: string; // Datetime.now()
    join_client_id: string; // 69
    name: string; // "Alt-Rock"
    now_playing: Song;
    num_users: number; // 2
    id: number; // 12
}

class TimeSync {
    type: string; // "SYNC"
    channel_id: number; // 12
    song_location: number; // 8621
    server_time: number; // 127606036519
}

class ErrorMessage {
    type: string; // "ERROR"
    caused_by_type: string; // "SONG" || "CLIENT" || "MESSAGE" || "CHANNEL" || "SYNC"
    message: string;
}