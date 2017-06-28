package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SessionPayload []map[string]string

func manageWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	// Initialize WebSocketConnection client with authentication
	wsConn := new(WebSocketConnection)
	wsConn.Init(r, ws)

	// register a client with the socket sender (possible message recipient)
	socketSender := new(SocketSender)
	socketSender.register(wsConn)

	log.Printf("webSocketConnections: %v", webSocketConnections)
	// listen on socket
	for {
		var message Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Printf("socket error: %v", err)
			delete(webSocketConnections, wsConn)
			break
		}

		// process incoming messages
		processWebSocketMessage(message, wsConn)
	}
}

func processWebSocketMessage(message Message, wsConn *WebSocketConnection) {
	messageType := message.Type
	messageDebug := false

	// @todo add configuration setting for message debug
	switch messageType {
	case "ping":
		break

	case "subscribe":
		var parameters MessageParameters
		json.Unmarshal([]byte(message.Parameters), &parameters)
		for _, subscription := range parameters.Subscriptions {
			folderPath := subscription.Folder
			if len(folderPath) == 0 {
				folderPath = subscription.FolderID
			}

			folderID, _ := strconv.Atoi(folderPath)
			if folderPath == "self" || folderPath == "[self]" {
				folderID = wsConn.ControllerID
			} else if strings.Contains(folderPath, "strip") { // @todo how to do hasattr in GO?  what does this check for exactly?
				resource := findResource(folderPath)
				folderID = resource.ID
			}

			if wsConn.hasAccess(folderID) {
				if messageDebug {
					log.Printf("subscribe folder: %s (%d)", folderPath, folderID)
				}
				wsConn.Subscriptions[folderID] = true
			}
		}
		break
	default:
		var folderID int
		if message.Folder != "" {
			folderName := message.Folder
			if messageDebug {
				log.Printf("message to folder: %v", folderName)
			}
			if strings.HasPrefix(folderName, "/") {
				if messageDebug {
					log.Printf("message to folder name: %v", folderName)
				}
				folder := findResource(folderName)
				if folder.ID > 0 {
					folderID = folder.ID
					if messageDebug {
						log.Printf("message to folder id: %v", folderID)
					}
				} else {
					log.Printf("message to unknown folder: (%v)", folderName)
				}
			}
		} else if wsConn.ControllerID > 0 {
			folderID = wsConn.ControllerID
		} else {
			log.Printf("message (%v) without folder or controller; discarding", messageType)
		}

		if wsConn.hasAccess(folderID) {
			timestamp := time.Now().Format(time.RFC850)
			mq := new(MessageQueue)
			// @todo does this need to be concurrent? It's a single insert.  Concurrent in Python / Flask app using a gevent for message_queue.add
			go mq.add(folderID, messageType, message.Parameters, wsConn.ControllerID, wsConn.UserID, timestamp)
		}

		break
	}
}
