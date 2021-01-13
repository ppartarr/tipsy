package users

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mcnijman/go-emailaddress"
	"github.com/ppartarr/tipsy/web/session"
)

// Login allows a user to login to their account
func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) error {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	expectedValues := []string{"email", "password"}
	if !formIsValid(r.Form, expectedValues) {
		return errors.New("you must submit a valid form")
	}

	// get user from email
	user, err := userService.getUser(r.Form["email"][0])
	if err != nil {
		return errors.New("could not get user " + r.Form["email"][0] + ": " + err.Error())
	}

	// validate email
	log.Println("validating email address")
	_, err = emailaddress.Parse(r.Form["email"][0])
	if err != nil {
		return errors.New("invalid email: %q" + r.Form["email"][0] + ": " + err.Error())
	}

	// TODO plug-in checkers here
	// check password
	log.Println("checking password")
	if !CheckPasswordHash(r.Form["password"][0], user.PasswordHash) {

		// increment login attempts
		if user.LoginAttempts < userService.config.RateLimit {
			userService.incrementLoginAttempts(user)
		} else {
			log.Println("user has attempted to login too many times")
			go func() {
				userService.PasswordRecovery(w, r)
			}()
			log.Println("this should return...")
			return errors.New("the limit of login attempts has been reached, please reset your password via the mail provided in the link " + r.Form["email"][0])
		}

		return errors.New("wrong password for user " + r.Form["email"][0])
	}

	// init session
	_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
	if err != nil {
		return errors.New("could not create a session for user " + r.Form["email"][0] + ": " + err.Error())
	}

	return nil
}
