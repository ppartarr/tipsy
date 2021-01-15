package users

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mcnijman/go-emailaddress"
	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/web/session"
)

// LoginForm represents a login form
type LoginForm struct {
	Email    string
	Password string
	Errors   map[string]string
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// Validate checks that the fields in the login form are set
func (form *LoginForm) Validate() bool {
	form.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Login"] = "Username and password incorrect"
	}

	if strings.TrimSpace(form.Password) == "" {
		form.Errors["Login"] = "Username and password incorrect"
	}

	return len(form.Errors) == 0
}

// Login allows a user to login to their account
func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) (form *LoginForm, err error) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return nil, errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	form = &LoginForm{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	if form.Validate() == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}
	// expectedValues := []string{"email", "password"}
	// if !formIsValid(r.Form, expectedValues) {
	// 	return errors.New("you must submit a valid form")
	// }

	// get user from email
	user, err := userService.getUser(r.Form["email"][0])
	if err != nil {
		return form, errors.New("could not get user " + r.Form["email"][0] + ": " + err.Error())
	}

	// validate email
	log.Println("validating email address")
	_, err = emailaddress.Parse(r.Form["email"][0])
	if err != nil {
		return form, errors.New("invalid email: %q" + r.Form["email"][0] + ": " + err.Error())
	}

	// use one of always, blacklist, optimal
	if userService.config.Checker.Always {
		log.Println("using always checker")
		// check password
		if !checkers.CheckAlways(r.Form["password"][0], user.PasswordHash) {

			// increment login attempts
			err = userService.checkLoginAttempts(w, r, user)
			if err != nil {
				return form, errors.New("failed to check login attemps")
			}
		}
	} else if userService.config.Checker.Blacklist != nil {
		log.Println("using blacklist checker")
		blacklist := checkers.LoadBlacklist(userService.config.Checker.Blacklist.File)

		// if password check fails, increment login attempts
		if !checkers.CheckBlacklist(r.Form["password"][0], user.PasswordHash, blacklist) {

			// increment login attempts
			err = userService.checkLoginAttempts(w, r, user)
			if err != nil {
				return form, errors.New("failed to check login attemps")
			}
		}
	} else if userService.config.Checker.Optimal != nil {
		log.Println("using optimal checker")
		frequencyBlacklist := checkers.LoadFrequencyBlackList(userService.config.Checker.Optimal.File)

		// if password check fails, increment login attempts
		if !checkers.CheckOptimal(r.Form["password"][0], user.PasswordHash, frequencyBlacklist, userService.config.Checker.Optimal.QthMostProbablePassword) {
			// increment login attempts
			err = userService.checkLoginAttempts(w, r, user)
			if err != nil {
				return form, errors.New("failed to check login attemps")
			}
		}
	}

	// init session
	_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
	if err != nil {
		return form, errors.New("could not create a session for user " + r.Form["email"][0] + ": " + err.Error())
	}

	return form, nil
}

func (userService *UserService) checkLoginAttempts(w http.ResponseWriter, r *http.Request, user *User) error {
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
