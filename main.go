package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("host", "0.0.0.0:8080", "Gotify server address")
var clientToken = flag.String("token", "", "Client token obtained from Gotify")

func main() {
	flag.Parse()
	log.SetFlags(0)

	wsURL := url.URL{
		Scheme: "ws",
		Host:   *addr,
		Path:   "/stream",
	}
	log.Printf("Websocket: Connecting to %s ...", wsURL.String())

	// TODO: GetAppIDs()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), http.Header{"X-Gotify-Key": []string{*clientToken}})
	if err != nil {
		log.Fatal("Websocket: failed to connect:", err)
	}
	defer c.Close()

	log.Printf("Websocket: connected!")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			ParseGotifyNotification(c)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupted")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Websocket: write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
