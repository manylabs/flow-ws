package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var webSocketConnections = make(map[*WebSocketConnection]bool)

var debugMode, _ = strconv.ParseBool(GetConfigVar("debugMode"))

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func initializeDatabase() *gorp.DbMap {
	// @todo sql strings in secrets file
	db, _ := sql.Open("sqlite3", GetConfigVar("SQLLitePath"))
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	return dbmap
}

// Message Define our message object
type Message struct {
	ID                 int             `json:"id"`
	Timestamp          string          `json:"timestamp"`
	SenderControllerID int             `json:"sender_controller_id"` //# the ID of the controller that created the message (if it was not created by a human/browser)
	SenderUserID       int             `json:"sender_user_id"`
	FolderID           int             `json:"folder_id"`
	Type               string          `json:"type"`
	Parameters         json.RawMessage `json:"parameters"`
	Folder             string          `json:"folder"`
}

type MessageParameters struct {
	AuthCode      string                          `json:"authCode"`
	Name          string                          `json:"name"`
	Folder        string                          `json:"folder"`
	FolderID      string                          `json:"folder_id"`
	Subscriptions []MessageParametersSubscription `json:"subscriptions"`
}

type MessageParametersSubscription struct {
	Folder   string `json:"folder"`
	FolderID string `json:"folder_id"`
}

type Key struct {
	ID                   int    `json:"id"`
	OrganizationID       string `json:"organization_id"`
	CreationUserID       string `json:"creation_user_id"`
	RevocationUserID     string `json:"revocation_user_id"`
	EmailAddress         string `json:"email_address"`
	AccessAsUserID       string `json:"access_as_user_id"`
	AccessAsControllerID string `json:"access_as_controller_id"`
	CreationTimestamp    string `json:"creation_timestamp"`
	RevocationTimestamp  string `json:"revocation_timestamp"`
	KeyPart              string `json:"key_part"`
	keyHash              string `json:"key_hash"`
	keyStorage           string `json:"key_storage"`
	keyStorageNonce      string `json:"key_storage_nonce"`
}

type WebSocketConnection struct {
	WS            *websocket.Conn `json:"ws"`
	UserID        int             `json:"id"`
	ControllerID  int             `json:"controller_id"`
	AuthMethod    int             `json:"auth_method"`
	Connected     int             `json:"connected"`
	Subscriptions map[string]bool `json:"subscription"`
}

func (wsConn *WebSocketConnection) Init(ws *websocket.Conn) {
	// @todo populate appropriately from authenticated user data
	wsConn.WS = ws
	wsConn.UserID = 4
}

func (wsConn *WebSocketConnection) hasAccess(FolderID string) bool {
	access := false
	if wsConn.ControllerID > 0 {
		// @todo fill in logic
	} else {
		// @todo fill in logic
	}
	return access
}

// Scan - Implement the database/sql scanner interface
func (msg *Message) Scan(value interface{}) error {
	log.Printf("%v", value)
	return nil
}

func main() {
	// // Configure websocket route
	http.HandleFunc("/api/v1/websocket", manageWebSocket)

	// // Start listening for incoming chat messages
	// go handleMessages()

	// go setUpSocketSender()

	listenPort := GetConfigVar("listenPort")
	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on", listenPort)
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func findKeyFromRequest(r *http.Request) Key {
	// @todo get user data from session cookie.  Session set in Flask app so logic needs porting.  Alternatively use a different method such as token
	// sessionCookie := ""
	// for _, cookie := range r.Cookies() {
	// 	if cookie.Name == "session" {
	// 		sessionCookie = cookie.Value
	// 	}
	// }

	return Key{}
}
