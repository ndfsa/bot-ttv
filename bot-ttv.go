package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

func main() {
	// get .env into environment
	godotenv.Load()


	// setup oauth flow
	tokenChan := make(chan string)
	setupAuth(tokenChan)

	// wait for access token
	accessToken := <-tokenChan

	// twitch bot logic
	prefix := os.Getenv("TWITCH_PREFIX")
	channelToJoin := os.Getenv("TWITCH_CHANNEL")

	client := twitch.NewClient("calc-bot", "oauth:"+accessToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("%s: %s\n", message.User.Name, message.Message)

		if !strings.HasPrefix(message.Message, prefix) {
			return
		}

		client.Say(message.Channel, strings.ToUpper(message.Message[7:]))
	})

	client.Join(channelToJoin)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func setupAuth(ch chan string) {
	// get secrets from environment
	clientId := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	// give url to the user
	fmt.Println("visit the following url to authorize")
	fmt.Printf("https://id.twitch.tv/oauth2/authorize?"+
		"response_type=code&"+
		"client_id=%s&"+
		"redirect_uri=http://localhost:3000/callback&"+
		"scope=chat%%3Aread+chat%%3Aedit\n", clientId)

	// setup http server to listen for responses from the authorization flow
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":3000",
		Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// get auth code from url
		authCode := r.URL.Query().Get("code")

		if authCode == "" {
			w.WriteHeader(http.StatusInternalServerError)
			panic("could not get auth code")
		}

		// get access token with auth code
		res, err := http.Post(fmt.Sprintf("https://id.twitch.tv/oauth2/token?"+
			"client_id=%s&"+
			"client_secret=%s&"+
			"code=%s&"+
			"grant_type=authorization_code&"+
			"redirect_uri=http://localhost:3000", clientId, clientSecret, authCode), "json", nil)

		// check response
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic("could not get token")
		}

		// parse body
		// TODO: parse body with a struct for the following object
		// {
		// 	"access_token": "rfx2uswqe8l4g1mkagrvg5tv0ks3",
		// 	"expires_in": 14124,
		// 	"refresh_token": "5b93chm6hdve3mycz05zfzatkfdenfspp1h1ar2xxdalen01",
		// 	"scope": [
		// 		"channel:moderate",
		// 		"chat:edit",
		// 		"chat:read"
		// 	],
		// 	"token_type": "bearer"
		// }
		var data map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic("could not read http response")
		}

		// send data
		ch <- data["access_token"].(string)
	})

	// final stage
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "authentication complete!")

		server.Shutdown(context.Background())
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
}
