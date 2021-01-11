package session

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var sessionMaxAge = int64(0)

const (
	keySessionExpires = "expires"
	keySignatureID    = "signatureID"
)

/*
 *	Public API - Setters
 */

// SetUserID sets the userID as a cookie on the response
// and persists it in the Store
func SetUserID(w http.ResponseWriter, r *http.Request, userID string) (string, error) {
	session, errStart := startSession(w, r, map[interface{}]interface{}{
		keyUserID: userID,
	})
	if errStart != nil {
		return "", errStart
	}

	Log.WithFields(logrus.Fields{
		"userID":    userID,
		"sessionID": session.ID,
	}).Info("setting userID on session")

	return session.ID, session.Save(r, w)
}

func startSession(w http.ResponseWriter, r *http.Request, values map[interface{}]interface{}) (*sessions.Session, error) {
	session := GetSession(r)
	values[keySessionExpires] = time.Now().Unix() + sessionMaxAge
	_, signatureExists := values[keySignatureID]

	if !signatureExists {
		values[keySignatureID] = randomID()
	}
	session.Values = values

	// generate a sessionID
	session.ID = randomID()

	return session, session.Save(r, w)
}

/*
 *	Public API - Getters
 */

// GetUserID returns the corresponding userID for a request
func GetUserID(w http.ResponseWriter, r *http.Request) string {
	return getKeyOrEmptyString(w, r, keyUserID)
}

// GetUserIDInt returns the corresponding userID for a request
func GetUserIDInt(w http.ResponseWriter, r *http.Request) (int, error) {
	// check for user sessions first
	rawID := GetUserID(w, r)
	if rawID == "" {
		return -1, errors.New("session missing or invalid")
	}

	userID, err := strconv.Atoi(rawID)
	if err != nil {
		return -1, errors.New("session missing or invalid")

	}

	return userID, nil
}

func GetSessionID(r *http.Request) string {
	s := GetSession(r)
	if s == nil {
		return ""
	}
	return s.ID
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	if session == nil || session.ID == "" {
		Log.WithFields(logrus.Fields{
			"url": r.URL,
		}).Warn("INVALID session")
	} else {
		Log.WithFields(logrus.Fields{
			"url": r.URL,
			"id":  session.ID,
		}).Info("VALID session")
	}
	return session
}

func GetNumSessions() int {
	files, err := ioutil.ReadDir(StorageDir)
	if err != nil {
		Log.WithError(err).Fatal("failed to read StorageDir")
	}

	return len(files)
}

/*
 *	Utils
 */

// Destroy a session
func Destroy(w http.ResponseWriter, r *http.Request) (err error) {
	session := GetSession(r)
	if session == nil {
		return errors.New("invalid session")
	}

	Log.Debug("DESTROY SESSION: " + session.ID)
	session.Values[keyUserID] = ""

	return session.Save(r, w)
}
