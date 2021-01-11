package session

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func randomID() string {
	return hex.EncodeToString(securecookie.GenerateRandomKey(32))
}

func getSignatureForSession(s *sessions.Session, subject ...string) (signature string, err error) {
	h := sha256.New()
	_, errWrite := h.Write([]byte(composeSubject(s, subject...)))
	if errWrite != nil {
		err = errors.New("could not hash")
		return
	}
	signature = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func getSignatureID(s *sessions.Session) string {
	signatureIDInterface, ok := s.Values[keySignatureID]
	if !ok {
		return ""
	}
	signatureID, okCast := signatureIDInterface.(string)
	if okCast {
		return signatureID
	}
	return ""
}

// compose salt + subjects + sessionID
func composeSubject(s *sessions.Session, subject ...string) string {
	return secret + strings.Join(subject, "-") + getSignatureID(s)
}

func getKeyOrEmptyString(w http.ResponseWriter, r *http.Request, key string) string {
	session := GetSession(r)
	if session != nil {
		v, ok := session.Values[key]
		if ok {
			return v.(string)
		}
	}
	return ""
}
