package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/gorcon/websocket"
)

const MockPasswordWebsocket = "password"

const MockCommandStatusResponseTextWebsocket = `status
hostname: Rust Server [DOCKER]
version : 2260 secure (secure mode enabled, connected to Steam3)
map     : Procedural Map
players : 0 (500 max) (0 queued) (0 joining)
id name ping connected addr owner violation kicks
`

func MockServerWebsocketHandlers() http.Handler {
	server := http.NewServeMux()

	var upgrader = gorilla.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	server.HandleFunc("/"+MockPasswordWebsocket, func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v\n", err)
			return
		}

		defer ws.Close()

		var response websocket.Message

		// Receive message.
		_, p, err := ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "gorilla: close 1006 (abnormal closure): unexpected EO") {
				log.Printf("read message error: %v\n", err)
			}
			return
		}

		var message websocket.Message
		if err := json.Unmarshal(p, &message); err != nil {
			// TODO: What Rust responses on read message fail?
			fmt.Println(string(p))
			log.Printf("unmarshal message error: %v\n", err)
			return
		}

		switch message.Message {
		case "status":
			response = websocket.Message{
				Message:    MockCommandStatusResponseTextWebsocket,
				Identifier: message.Identifier,
				Type:       "Generic",
			}
		case "deadline":
			time.Sleep(websocket.DefaultDeadline + 1*time.Second)
			response = websocket.Message{
				Message:    fmt.Sprintf("sleep for %d secends", websocket.DefaultDeadline+1*time.Second),
				Identifier: message.Identifier,
				Type:       "Generic",
			}
		default:
			response = websocket.Message{
				Message:    fmt.Sprintf("Command '%s' not found", message.Message),
				Identifier: message.Identifier,
				Type:       "Warning",
			}
		}

		js, err := json.Marshal(response)
		if err != nil {
			log.Printf("marshal response error: %v\n", err)
			return
		}

		if err := ws.WriteMessage(gorilla.TextMessage, js); err != nil {
			log.Printf("write response error: %v\n", err)
			return
		}
	})

	return server
}
