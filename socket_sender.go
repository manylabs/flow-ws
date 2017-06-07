package main

import (
	"encoding/json"
	"log"
	"time"
)

type SocketSender struct{}

//init
// @todo is this needed as we have global webSocketConnections

// register a client (possible message recipient)
func (s SocketSender) register(wsConn *WebSocketConnection) {
	log.Printf("new client registered")
	webSocketConnections[wsConn] = true
}

// send a message to a specific client (using websocket connection specified in wsConn)
func (s SocketSender) send(wsConn *WebSocketConnection, message Message) {
	err := wsConn.WS.WriteJSON(message)
	if err != nil {
		if wsConn.ControllerID > 0 {
			log.Printf("send error to controller ID: %v", wsConn.ControllerID)
		} else {
			log.Printf("send error")
		}
		wsConn.WS.Close()
		if _, ok := webSocketConnections[wsConn]; ok {
			delete(webSocketConnections, wsConn)
		}
	}
}

// send an error message back to a client
func (s SocketSender) sendError(wsConn *WebSocketConnection, messageText string) {
	message := Message{}

	parametersStruct := map[string]interface{}{
		"message": messageText,
	}
	parametersJSON, _ := json.Marshal(parametersStruct)

	message.Type = "error"
	// @todo use proper time format
	message.Timestamp = time.Now().Format(time.RFC850)
	message.Parameters = parametersJSON

	s.send(wsConn, message)
}

// this function sits in a loop, waiting for messages that need to be sent out to subscribers)
func (s SocketSender) run() {
	for {
		time.Sleep(1 * time.Second)
		mq := MessageQueue{}
		messages := mq.receive()
		for _, message := range messages {
			if debugMode == true {
				log.Printf("message type: %v, folder: %v", message.Type, message.FolderID)
			}
			for wsConn := range webSocketConnections {
				if ClientIsSubscribed(message, wsConn) {
					messageStruct := Message{}
					messageStruct.Type = message.Type
					// @todo use proper time format
					messageStruct.Timestamp = message.Timestamp
					messageStruct.Parameters = message.Parameters
					s.send(wsConn, messageStruct)
					if debugMode == true {
						if wsConn.ControllerID > 0 {
							log.Printf("sending message to controller; type: %v", message.Type)
						} else {
							log.Printf("sending message to browser; type: %v", message.Type)
						}
					}
				}
			}
		}
	}

}

// returns True if the given message should be sent to the given client (based on its current subscriptions)
func ClientIsSubscribed(message Message, wsConn *WebSocketConnection) bool {
	if message.SenderControllerID > 0 {
		if (wsConn.ControllerID > 0) && message.SenderControllerID == wsConn.ControllerID {
			return false
		}
		if (wsConn.UserID > 0) && message.SenderUserID == wsConn.UserID {
			return false
		}
	}
	if _, ok := wsConn.Subscriptions[string(message.FolderID)]; ok {
		return ok
	}
	return false
}

func setUpSocketSender() {
	s := SocketSender{}
	s.run()
}
