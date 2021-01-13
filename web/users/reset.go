package users

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// PasswordReset validates the token then updates the user's password
func (userService *UserService) PasswordReset(w http.ResponseWriter, r *http.Request) error {

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email", "password", "password-copy"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// check that password & password copy match
	if r.Form["password"][0] != r.Form["password-copy"][0] {
		return errors.New("password and password copy don't match")
	}

	// get token hash from form & get associated token from bucket
	tokenHash := r.Form["token"][0]
	token, err := userService.getToken(tokenHash)
	if err != nil {
		return errors.New("could not get token")
	}

	log.Println(token)

	// check that the token matches the submitted user
	if r.Form["email"][0] != token.Email {
		return errors.New("submitted email & token email don't match")
	}

	// check if the TTL is still valid
	if token.CreatedAt.Add(token.TTL).Before(time.Now().Local()) {
		return errors.New("the token expired at " + token.CreatedAt.Add(token.TTL).String())
	}

	// invalidate password reset token
	err = userService.deleteToken(token)
	if err != nil {
		return errors.New("error deleting the token" + err.Error())
	}

	// get user from db
	user, err := userService.getUser(token.Email)
	if err != nil {
		return errors.New("error getting the user: " + err.Error())
	}

	// hash the input password
	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		return errors.New("failed to hash the input password: " + err.Error())
	}

	// set new password & set login attempts to 0
	user.PasswordHash = passwordHash
	user.LoginAttempts = 0

	// update the user's password
	err = userService.updateUser(user)
	if err != nil {
		return errors.New("failed to update the user's password: " + err.Error())
	}

	log.Println("password has been reset")

	return nil
}
