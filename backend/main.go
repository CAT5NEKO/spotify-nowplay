package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type NowPlayingResponse struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	Album         string `json:"album"`
	Url           string `json:"url"`
	IsPlaying     bool   `json:"isPlaying"`
	AlbumCoverURL string `json:"albumCoverURL"`
}

func nowPlayingHandler(w http.ResponseWriter, r *http.Request) {

	isPlaying, title, artist, album, url, _, albumCoverURL := getSpotifyNP()

	response := NowPlayingResponse{
		Title:         title,
		Artist:        artist,
		Album:         album,
		Url:           url,
		IsPlaying:     isPlaying,
		AlbumCoverURL: albumCoverURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var authCompleted = make(chan struct{})
var stopSendingUpdates = make(chan struct{})

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan NowPlayingResponse)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		select {
		case msg := <-broadcast:
			sendDataToClients(msg)
		}
	}
}

func sendDataToClients(msg NowPlayingResponse) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func sendNowPlayingUpdatesToWebSocket() {
	for {
		select {
		case <-stopSendingUpdates:

			return
		default:
			isPlaying, title, artist, album, url, _, albumCoverURL := getSpotifyNP()
			response := NowPlayingResponse{
				Title:         title,
				Artist:        artist,
				Album:         album,
				Url:           url,
				IsPlaying:     isPlaying,
				AlbumCoverURL: albumCoverURL,
			}

			broadcast <- response

			time.Sleep(5 * time.Second)
		}
	}
}

func main() {

	http.HandleFunc("/now-playing", nowPlayingHandler)

	http.HandleFunc("/ws", handleConnections)

	go func() {
		http.HandleFunc("/login", spotify_login)
		auth_code := make(chan string)
		go pass_callback(auth_code)
		handleCallback := spotify_callback(auth_code)
		http.HandleFunc("/callback", handleCallback)

		err := http.ListenAndServe("0.0.0.0:4400", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	godotenv.Load(".env")
	if os.Getenv("SPOTIFY_CLIENT_ID") == "" || os.Getenv("SPOTIFY_CLIENT_SECRET") == "" {
		log.Fatal("Spotifyで必要な資格要件が不足しています。envを修正してください。")
	} else if os.Getenv("SPOTIFY_REFRESH_TOKEN") == "" {
		fmt.Println("`SPOTIFY_REFRESH_TOKEN` がセットされていません。以下よりセットしてください.")
		values := url.Values{}
		values.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
		values.Add("response_type", "code")
		values.Add("redirect_uri", "http://localhost:4400/callback")
		values.Add("scope", "user-read-playback-state user-read-currently-playing")
		fmt.Println("https://accounts.spotify.com/authorize?" + values.Encode())

		<-authCompleted
	} else {

		go sendNowPlayingUpdatesToWebSocket()
	}

	<-stopSendingUpdates
}
