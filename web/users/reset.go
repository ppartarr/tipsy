package users

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nbutton23/zxcvbn-go"
	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/correctors"
)

// ResetForm represents a reset form
type ResetForm struct {
	Email        string
	Password     string
	PasswordCopy string
	Errors       map[string]string
	Token        string
}

// Validate checks that the fields in the login form are set
func (form *ResetForm) Validate(config *config.Server, zxcvbnScore int, token *ResetToken) bool {
	form.Errors = make(map[string]string)

	// check that username is an email address
	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Email"] = "Email must be valid"
	}

	// check if the TTL is still valid
	if token.CreatedAt.Add(token.TTL).Before(time.Now().Local()) {
		form.Errors["Email"] = "the token expired at " + token.CreatedAt.Add(token.TTL).String()
	}

	if strings.TrimSpace(form.Password) == "" {
		form.Errors["Password"] = "Password cannot be empty"
	}

	// check that password & password copy match
	if form.Password != form.PasswordCopy {
		form.Errors["Password"] = "Passwords should match"
	}

	// check if password is in blacklist
	if config.Checker.Blacklist != nil {
		blacklist := checkers.LoadBlacklist(config.Checker.Blacklist.File)
		if correctors.StringInSlice(form.Password, blacklist) {
			form.Errors["Password"] = "Password is forbidden"
		}
	}

	// only register password if strength estimation is high enough
	// TODO use same zxcvbn in front-end and backend
	score := zxcvbn.PasswordStrength(form.Password, []string{form.Email})
	log.Println(score)
	if score.Score < zxcvbnScore {
		form.Errors["Password"] = "Password should be at least " + strconv.Itoa(zxcvbnScore) + "/4 in zxcvbn"
	}

	// check if submitted token & user match
	if token.Token != form.Token {
		log.Println("submitted token doesn't match with the user")
		form.Errors["Email"] = "Account hasn't requested a password reset"
	}

	return len(form.Errors) == 0
}

// PasswordReset validates the token then updates the user's password
func (userService *UserService) PasswordReset(w http.ResponseWriter, r *http.Request) (form *ResetForm, err error) {

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	form = &ResetForm{
		Email:        r.PostFormValue("email"),
		Password:     r.PostFormValue("password"),
		PasswordCopy: r.PostFormValue("password-copy"),
		Token:        r.PostFormValue("token"),
	}

	// get user from db using submitted email
	user, err := userService.getUser(form.Email)
	if err != nil {
		form.Errors = make(map[string]string)
		log.Println("could not get user" + form.Email + ": " + err.Error())
		form.Errors["Email"] = "Email does not exist"
		return form, errors.New("you must submit a valid form")
	}

	// validate form
	if form.Validate(userService.config, userService.config.Web.Register.Zxcvbn, user.ResetToken) == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	log.Println(user.ResetToken.Token)

	// hash the input password
	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		return nil, errors.New("failed to hash the input password: " + err.Error())
	}

	// delete token, set new password, set login attempts to 0
	user.ResetToken = nil
	user.PasswordHash = passwordHash
	user.LoginAttempts = 0

	// update the user
	err = userService.updateUser(user)
	if err != nil {
		return nil, errors.New("failed to update the user's password: " + err.Error())
	}

	log.Println("password has been reset")

	return nil, nil
}
