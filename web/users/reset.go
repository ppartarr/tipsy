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
func (form *ResetForm) Validate(blacklistFile string, zxcvbnScore int, token *Token) bool {
	form.Errors = make(map[string]string)

	// check that username is an email address
	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Email"] = "Email must be valid"
	}

	// check that the token matches the submitted user
	if form.Email != token.Email {
		form.Errors["Email"] = "submitted email & token email don't match"
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
	blacklist := checkers.LoadBlacklist(blacklistFile)
	if checkers.StringInSlice(form.Password, blacklist) {
		form.Errors["Password"] = "Password is forbidden"
	}

	// only register password if strength estimation is high enough
	// TODO use same zxcvbn in front-end and backend
	score := zxcvbn.PasswordStrength(form.Password, []string{form.Email})
	log.Println(score)
	if score.Score < zxcvbnScore {
		form.Errors["Password"] = "Password should be at least " + strconv.Itoa(zxcvbnScore) + "/4 in zxcvbn"
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

	// get token hash from form & get associated token from bucket
	tokenHash := form.Token
	token, err := userService.getToken(tokenHash)
	if err != nil {
		return nil, errors.New("could not get token")
	}

	// validate form
	if form.Validate(userService.config.Web.Register.Blacklist, userService.config.Web.Register.Zxcvbn, token) == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	log.Println(token)

	// invalidate password reset token
	err = userService.deleteToken(token)
	if err != nil {
		return nil, errors.New("error deleting the token" + err.Error())
	}

	// get user from db
	user, err := userService.getUser(token.Email)
	if err != nil {
		return nil, errors.New("error getting the user: " + err.Error())
	}

	// hash the input password
	passwordHash, err := HashPassword(r.Form["password"][0])
	if err != nil {
		return nil, errors.New("failed to hash the input password: " + err.Error())
	}

	// set new password & set login attempts to 0
	user.PasswordHash = passwordHash
	user.LoginAttempts = 0

	// update the user's password
	err = userService.updateUser(user)
	if err != nil {
		return nil, errors.New("failed to update the user's password: " + err.Error())
	}

	log.Println("password has been reset")

	return nil, nil
}
