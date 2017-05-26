package main

import "log"

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
	_, err := dbmap.Select(&messages, "select * from messages order by id limit 0, 1")
	if err != nil {
		log.Printf("err: %v", err)
	}
	// for _, msg := range messages {
	// 	log.Printf("msg: %v", msg)

	// }
	return messages
}

func (m MessageQueue) add(folderID string, messageType string, parameters string, senderControllerID string, senderUserID string, timestamp string) {
	// @todo persist to db
}
