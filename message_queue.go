package main

import (
	"encoding/json"
	"log"
)

type MessageQueue struct {
	lastMessageID  int
	cleanUpRunning int
	startTimeStamp string // @todo change this to time object
	wakeEvent      int
}

func (m MessageQueue) receive() []Message {
	//get messages
	// initialize the DbMap
	dbmap := initializeDatabase()
	defer dbmap.Db.Close()
	var messages []Message

	if m.lastMessageID > 0 {
		_, err := dbmap.Select(&messages, "select * from messages WHERE id > :last_message_id ORDER BY id", map[string]interface{}{
			"last_message_id": m.lastMessageID,
		})
		if err != nil {
			log.Printf("err: %v", err)
		}
	} else {
		_, err := dbmap.Select(&messages, "select * from messages WHERE id > :start_timestamp ORDER BY id", map[string]interface{}{
			"start_timestamp": m.startTimeStamp,
		})
		if err != nil {
			log.Printf("err: %v", err)
		}
	}
	if len(messages) > 0 {
		m.lastMessageID = messages[len(messages)-1].ID
	}
	return messages
}

func (m MessageQueue) add(folderID int, messageType string, parameters string, senderControllerID int, senderUserID int, timestamp string) {
	// @todo persist to db
	message := Message{}
	message.SenderControllerID = senderControllerID
	message.SenderUserID = senderUserID
	message.Timestamp = timestamp
	message.Type = messageType
	message.FolderID = folderID
	message.Parameters, _ = json.Marshal(parameters)
	// @todo add parameters
	// @tood if no timestamp
	dbmap := initializeDatabase()
	defer dbmap.Db.Close()
	err := dbmap.Insert(&message)
	if err != nil {
		log.Printf("MessageQueue.add err: %v", err)
	}
}
