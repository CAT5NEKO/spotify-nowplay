package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

func spotify_login(w http.ResponseWriter, req *http.Request) {
	values := url.Values{}
	values.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	values.Add("response_type", "code")
	values.Add("redirect_uri", "http://localhost:4400/callback")
	values.Add("scope", "user-read-playback-state user-read-currently-playing")

	http.Redirect(w, req, "https://accounts.spotify.com/authorize?"+values.Encode(), http.StatusFound)
}

func spotify_callback(auth_code chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query().Get("code")
		auth_code <- query

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=utf8")

		w.Write([]byte("処理が完了しました。この画面を閉じることができます。再起動してください。"))
	}
}

func pass_callback(auth_code chan string) {
	for item := range auth_code {
		save_refresh_token(item)
	}
}
func save_refresh_token(auth_code string) {

	values := make(url.Values)
	values.Set("grant_type", "authorization_code")
	values.Set("code", auth_code)

	values.Set("redirect_uri", "http://localhost:4400/callback")
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatalf("POSTリクエストの送信に失敗しました。: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))))))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("トークン変換リクエストに失敗しました。: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンスボディの読み取りに失敗しました。: %s", err)
	}

	var jsonObj interface{}
	if err := json.Unmarshal(body, &jsonObj); err != nil {
		fmt.Println(string(body))
		log.Fatalf("JSONボディにパースする所で問題が発生しました。: %s\nResponse body: %s", err, string(body))
	}

	refresh_token := jsonObj.(map[string]interface{})["refresh_token"].(string)
	refresh_token_env, err := godotenv.Unmarshal(fmt.Sprintf("SPOTIFY_CLIENT_ID=%s\nSPOTIFY_CLIENT_SECRET=%s\nSPOTIFY_REFRESH_TOKEN=%s\n", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), refresh_token))

	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Write(refresh_token_env, "./.env")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func get_spotify_access_token() string {
	values := make(url.Values)
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatalf("POSTリクエストが送信できませんでした。: %s", err)
	}

	spotify_auth_string := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", spotify_auth_string))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jsonObj interface{}
	if err := json.Unmarshal(body, &jsonObj); err != nil {
		fmt.Println(string(body))
		log.Fatal(err)
	}

	if isNil(jsonObj.(map[string]interface{})["access_token"]) {
		fmt.Println(body)
		os.Exit(1)
	}

	return jsonObj.(map[string]interface{})["access_token"].(string)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func get_spotify_np() (is_playing bool, title string, artist string, album string, url string, progress float64) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		log.Fatalf("HTTPリクエストの作成に失敗しました。: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", get_spotify_access_token()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("HTTPリクエストで問題が発生しました。: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if len(body) == 0 {
		log.Println("Empty response from Spotify API")
		return false, "", "", "", "", 0
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Error: オーソライズに失敗しています。`SPOTIFY_REFRESH_TOKEN` を確認してください。")
	}

	var jsonObj interface{}
	if err := json.Unmarshal(body, &jsonObj); err != nil {

		fmt.Println(string(body))
		log.Fatalf("JSON unmarshal で問題が生じました。: %s\nResponse body: %s", err, string(body))
	}

	if isNil(jsonObj) || isNil(jsonObj.(map[string]interface{})["is_playing"]) {

		fmt.Println(string(body))
		log.Println("エラーが発生しました。Spotify が再生中でない可能性があります。")
		return false, "", "", "", "", 0
	}

	is_playing = jsonObj.(map[string]interface{})["is_playing"].(bool)

	if is_playing {
		title = jsonObj.(map[string]interface{})["item"].(map[string]interface{})["name"].(string)

		artists := jsonObj.(map[string]interface{})["item"].(map[string]interface{})["artists"].([]interface{})
		artistList := make([]string, len(artists))
		for i, artist := range artists {
			artistList[i] = artist.(map[string]interface{})["name"].(string)
		}
		artist = strings.Join(artistList, ", ")

		album = jsonObj.(map[string]interface{})["item"].(map[string]interface{})["album"].(map[string]interface{})["name"].(string)

		url = jsonObj.(map[string]interface{})["item"].(map[string]interface{})["external_urls"].(map[string]interface{})["spotify"].(string)

		progress = jsonObj.(map[string]interface{})["progress_ms"].(float64)
	} else {
		is_playing = false

		title, artist, album = "", "", ""
	}

	return is_playing, title, artist, album, url, progress
}
