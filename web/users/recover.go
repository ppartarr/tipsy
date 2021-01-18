package users

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ppartarr/tipsy/mail"
)

// PasswordRecovery sends the user an email containing a password reset link
func (userService *UserService) PasswordRecovery(w http.ResponseWriter, r *http.Request) error {
	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// get user from db using submitted email
	user, err := userService.getUser(r.Form["email"][0])
	if err != nil {
		return errors.New("could not get user" + r.Form["email"][0] + ": " + err.Error())
	}

	log.Println("there")

	log.Println(userService.config.TokenValidity)
	log.Println(userService.config.TokenValidity * time.Minute)

	token := &Token{
		Email:     r.Form["email"][0],
		TTL:       userService.config.TokenValidity,
		CreatedAt: time.Now().Local(),
	}

	// generate token
	token.Token = generateToken(user.Email)

	// store the token
	log.Println(token)
	err = userService.storeToken(token)
	if err != nil {
		return errors.New("could not store token for user " + r.Form["email"][0] + ": " + err.Error())
	}

	// send password reset mail
	go func() {
		log.Println("sending password reset email")
		// i hope this is safe
		url := "http://localhost:8000/reset.html?token=" + token.Token
		mail.Send(user.Email, "Tipsy password reset", mail.GeneratePasswordResetMail(url))
	}()

	return nil
}

// checks that the form isn't empty...
// TODO improve this by adding length requirements
func formIsValid(form url.Values, expectedValues []string) bool {
	for _, value := range expectedValues {
		if len(form) == 0 || len(form[value]) == 0 {
			return false
		}
	}

	return true
}
