package users

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ppartarr/tipsy/mail"
)

// RecoveryForm is a recovery form
type RecoveryForm struct {
	Email  string
	Errors map[string]string
}

// Validate checks that the fields in the login form are set
func (form *RecoveryForm) Validate() bool {
	form.Errors = make(map[string]string)

	// check that username is an email address
	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Email"] = "Email must be valid"
	}

	return len(form.Errors) == 0
}

// PasswordRecovery sends the user an email containing a password reset link
func (userService *UserService) PasswordRecovery(w http.ResponseWriter, r *http.Request) (form *RecoveryForm, err error) {
	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	form = &RecoveryForm{
		Email: r.PostFormValue("email"),
	}
	if form.Validate() == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	// get user from db using submitted email
	user, err := userService.getUser(form.Email)
	if err != nil {
		log.Println("could not get user" + form.Email + ": " + err.Error())
		form.Errors["Email"] = "Email does not exist"
		return form, errors.New("you must submit a valid form")
	}

	log.Println("there")

	log.Println(userService.config.Web.Reset.TokenValidity)
	log.Println(userService.config.Web.Reset.TokenValidity * time.Minute)

	token := &Token{
		Email:     form.Email,
		TTL:       userService.config.Web.Reset.TokenValidity,
		CreatedAt: time.Now().Local(),
	}

	// generate token
	token.Token = generateToken(user.Email)

	// store the token
	log.Println(token)
	err = userService.storeToken(token)
	if err != nil {
		return nil, errors.New("could not store token for user " + form.Email + ": " + err.Error())
	}

	// send password reset mail
	go func() {
		log.Println("sending password reset email")
		// i hope this is safe
		url := "http://localhost:8000/reset.html?token=" + token.Token
		mail.Send(user.Email, "Tipsy password reset", mail.GeneratePasswordResetMail(url))
	}()

	return nil, nil
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
