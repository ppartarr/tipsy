package users

import (
	"errors"
	"log"
	"net/http"
)

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) error {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return errors.New("must use POST http method")
	}

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

	// check that email isn't already registered
	user, err := userService.getUser(r.Form["email"][0])
	if user != nil {
		return errors.New("user already with email " + r.Form["email"][0] + " is already registered")
	}

	// TODO only register password if it is 3/4 in zxcbn-go

	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		log.Println(err.Error())
	}

	// create new user from request then save in db
	user = &User{
		Email:         r.Form["email"][0],
		PasswordHash:  passwordHash,
		LoginAttempts: 0,
	}

	userService.createUser(user)

	return nil
}
