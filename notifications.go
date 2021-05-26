package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/haklop/gnotifier"
)

// TODO: Import these from gotify/server
type GotifyMessage struct {
	AppID    int    `json:"appid"`
	Date     string `json:"date"`
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	Title    string `json:"title"`
}

type GotifyApplication struct {
	ID          int    `json:"id"`
	AppToken    string `json:"token"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Internal    bool   `json:"internal"`
	Image       string `json:"image"`
}

func sendNotification(title string, message string) {
	notification := gnotifier.Notification(title, message)
	notification.GetConfig().Expiration = 2000
	notification.GetConfig().ApplicationName = "Gotify"
	notification.Push()
}

func GetAppIDs() []GotifyApplication {
	appEndpoint := fmt.Sprintf("http://%s/application?token=%s", *addr, *clientToken)

	res, err := http.Get(appEndpoint)
	if err != nil {
		log.Print("Could not retrieve list of Gotify applications:", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Invalid response for application list:", err)
	}

	var apps []GotifyApplication
	if err := json.Unmarshal(body, &apps); err != nil {
		panic("GetAppIDs: Failed to decode JSON")
	}
	return apps
}

func ParseGotifyNotification(c *websocket.Conn) {
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("Websocket: read error:", err)
		return
	}

	var msg GotifyMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		panic("GotifyMessage: Failed to decode JSON")
	}
	sendNotification(msg.Title, msg.Message)
}
