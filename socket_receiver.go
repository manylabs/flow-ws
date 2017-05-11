package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func manageWebSocket(w http.ResponseWriter, r *http.Request) {
	// @todo check authentication in request
	// get session info

	// find key based on passed session info
	key := findKeyFromRequest(r)
	log.Printf("cookie: %v", key)

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	// Initialize WebSocketConnection client
	wsConn := new(WebSocketConnection)
	wsConn.Init(ws)

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

		// // Send the newly received message to the broadcast channel
		// broadcast <- msg
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

			folderID := folderPath
			if folderPath == "self" || folderPath == "[self]" {
				folderID = string(wsConn.ControllerID)
			} else if strings.Contains(folderPath, "strip") { // @todo how to do hasattr in GO?  what does this check for exactly?
				// @todo implement find_resource
				// resource := findResource(folderID)
				// folderId = resource.ID
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
		folderID := ""
		if message.Folder != "" {
			folderName := message.Folder
			if messageDebug {
				log.Printf("message to folder: %v", folderName)
			}
			if strings.HasPrefix(folderName, "/") {
				// @todo fill in logic
			}
		} else if wsConn.ControllerID > 0 {
			folderID := wsConn.ControllerID
			log.Printf("folderID = ControllerID", folderID)
		} else {
			log.Printf("message (%v) without folder or controller; discarding", messageType)
		}

		if wsConn.hasAccess(folderID) {
			var parameters MessageParameters
			json.Unmarshal([]byte(message.Parameters), &parameters)

			// add to message queue
			// mq := new(MessageQueue)
			// mq.add(folderID, messageType, parameters)
		}

		break
	}
}
