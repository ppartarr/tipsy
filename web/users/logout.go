package users

import (
	"errors"
	"net/http"

	"github.com/ppartarr/tipsy/web/session"
)

// Logout logs a user out
func (userService *UserService) Logout(w http.ResponseWriter, r *http.Request) error {
	// check if there is a current session
	if session.GetSessionID(r) == "" {
		// destroy the existing session
		err := session.Destroy(w, r)
		if err != nil {
			return errors.New("could not destroy session")
		}
	}

	return nil
}
