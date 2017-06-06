package main

import (
	"bytes"
	"compress/zlib"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// NullTime is taken from sql/pq
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// @todo handle errors in utils
func getSessionCookieFromRequest(r *http.Request) map[string]interface{} {

	sessionPayloadMap := make(map[string]interface{})
	sessionCookie := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session" {
			sessionCookie = cookie.Value
			sessionPayloadMap = decodeSessionCookie(sessionCookie)
		}
	}

	return sessionPayloadMap
}

func decodeSessionCookie(sessionCookie string) map[string]interface{} {
	parts := strings.Split(sessionCookie, ".")

	str := parts[1]
	str += "==="

	decoded, _ := base64.RawURLEncoding.DecodeString(str)

	reader := bytes.NewReader(decoded)
	zlibReader, _ := zlib.NewReader(reader)
	jsonPayloadBytes, _ := ioutil.ReadAll(zlibReader)

	sessionPayloadMap := make(map[string]interface{})

	_ = json.Unmarshal(jsonPayloadBytes, &sessionPayloadMap)

	return sessionPayloadMap
}

func findKey(keyText string) Key {
	keyPart := keyText[:3] + keyText[len(keyText)-3:]
	dbmap := initializeDatabase()
	defer dbmap.Db.Close()
	var key Key
	err := dbmap.SelectOne(&key, "select * from keys where key_part=:key_part and revocation_timestamp IS NULL", map[string]interface{}{
		"key_part": keyPart,
	})
	if err != nil {
		log.Printf("findKey err: %v", err)
	}
	// @todo check key authentication and only return key if true
	return key
}

func findResource(fileName string) Resource {
	fileName = strings.Trim(fileName, "/")
	parts := strings.Split(fileName, "/")
	parent := Resource{}

	dbmap := initializeDatabase()
	defer dbmap.Db.Close()

	resource := Resource{}
	for part := range parts {
		if (parent != Resource{}) {
			err := dbmap.SelectOne(&resource, "select * from resource where parent_id=:parent_id AND name=:=name AND deleted = false", map[string]interface{}{
				"parent_id": parent.ID,
				"name":      part,
			})
			if err != nil {
				log.Printf("findResource err: %v", err)
			}
		} else {
			err := dbmap.SelectOne(&resource, "select * from resource where parent_id=:parent_id AND name=:=name AND deleted = false", map[string]interface{}{
				"parent_id": nil,
				"name":      part,
			})
			if err != nil {
				log.Printf("findResource err: %v", err)
			}
		}
		// need to handle when no results found
		// need to handle when more than 1 resource.. and report multiple results found.
		parent = resource
	}
	return parent
}
