package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type NowPlayingResponse struct {
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	Album         string `json:"album"`
	Url           string `json:"url"`
	IsPlaying     bool   `json:"isPlaying"`
	AlbumCoverURL string `json:"albumCoverURL"`
}

var (
	authCompleted       = make(chan struct{})
	stopSendingUpdates  = make(chan struct{})
	upgrader            = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients             = make(map[*websocket.Conn]bool)
	clientsMutex        sync.Mutex
	broadcast           = make(chan NowPlayingResponse)
)

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

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocketアップグレードエラー:", err)
		return
	}

	clientsMutex.Lock()
	clients[ws] = true
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, ws)
		clientsMutex.Unlock()
		ws.Close()
	}()

	for {
		if _, _, err := ws.NextReader(); err != nil {
			break
		}
	}
}

func sendDataToClients(msg NowPlayingResponse) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			fmt.Printf("WebSocket送信エラー: %v\n", err)
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
	_ = godotenv.Load(".env")

	http.HandleFunc("/now-playing", nowPlayingHandler)
	http.HandleFunc("/ws", handleConnections)

	go func() {
		for msg := range broadcast {
			sendDataToClients(msg)
		}
	}()

	if os.Getenv("SPOTIFY_CLIENT_ID") == "" || os.Getenv("SPOTIFY_CLIENT_SECRET") == "" {
		log.Fatal("SpotifyのクライアントID/シークレットが不足しています。.envを確認してください。")
	} else if os.Getenv("SPOTIFY_REFRESH_TOKEN") == "" {
		fmt.Println("`SPOTIFY_REFRESH_TOKEN` が設定されていません。以下のURLから認証してください：")
		values := url.Values{}
		values.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
		values.Add("response_type", "code")
		values.Add("redirect_uri", "http://127.0.0.1:4400/callback")
		values.Add("scope", "user-read-playback-state user-read-currently-playing")
		fmt.Println("https://accounts.spotify.com/authorize?" + values.Encode())
		auth_code := make(chan string)
		go pass_callback(auth_code)
		handleCallback := spotify_callback(auth_code)
		http.HandleFunc("/callback", handleCallback)
		http.HandleFunc("/login", spotify_login)

		go func() {
			err := http.ListenAndServe("0.0.0.0:4400", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()

		<-authCompleted
	} else {
		go sendNowPlayingUpdatesToWebSocket()
		go func() {
			err := http.ListenAndServe("0.0.0.0:4400", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	<-stopSendingUpdates
}
