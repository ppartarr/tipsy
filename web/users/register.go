package users

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/nbutton23/zxcvbn-go"
	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/checkers/typtop"
)

// RegistrationForm represents a login form
type RegistrationForm struct {
	Email        string
	Password     string
	PasswordCopy string
	Errors       map[string]string
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// Validate checks that the fields in the login form are set
func (form *RegistrationForm) Validate(blacklistFile string, zxcvbnScore int) bool {
	form.Errors = make(map[string]string)

	// check that username is an email address
	match := rxEmail.Match([]byte(form.Email))
	if match == false {
		form.Errors["Email"] = "Email must be valid"
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

// Register allows a user to register a new account
func (userService *UserService) Register(w http.ResponseWriter, r *http.Request) (form *RegistrationForm, err error) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return nil, errors.New("must use POST http method")
	}

	// check that a valid form was submitted
	r.ParseForm()
	log.Println(r.Form)
	form = &RegistrationForm{
		Email:        r.PostFormValue("email"),
		Password:     r.PostFormValue("password"),
		PasswordCopy: r.PostFormValue("password-copy"),
	}

	// validate form
	if form.Validate(userService.config.Web.Register.Blacklist, userService.config.Web.Register.Zxcvbn) == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	if userService.config.Checker.TypTop != nil {
		log.Println("running typtop mode")

		// check that email isn't already registered
		typtopUser, err := userService.getTypTopUser(form.Email)
		if typtopUser != nil {
			log.Println("typtop user already with email " + form.Email + " is already registered")
			form.Errors["Email"] = "Email already registered"
			return form, errors.New("you must submit a valid form")
		}

		// init the checker service
		log.Println("init checker service")
		Checker := typtop.NewChecker(userService.config.Checker.TypTop, userService.config.Typos)

		// register the password for typtop
		typtopState, privateKey := Checker.Register(form.Password)

		// create new user from request then save in db
		typtopUser = &typtop.User{
			Email:         form.Email,
			LoginAttempts: 0,
			State:         typtopState,
			PrivateKey:    privateKey,
		}

		err = userService.createTypTopUser(typtopUser)
		if err != nil {
			return nil, errors.New("couldn't create typtop user")
		}
	} else {

		// check that email isn't already registered
		user, err := userService.getUser(form.Email)
		if user != nil {
			log.Println("user already with email " + form.Email + " is already registered")
			form.Errors["Email"] = "Email already registered"
			return form, errors.New("you must submit a valid form")
		}

		passwordHash, err := HashPassword(form.Password)
		if err != nil {
			return nil, errors.New("couldn't hash password: " + err.Error())
		}

		// create new user from request then save in db
		user = &User{
			Email:         form.Email,
			PasswordHash:  passwordHash,
			LoginAttempts: 0,
		}

		err = userService.createUser(user)
		if err != nil {
			return nil, errors.New("couldn't create user")
		}
	}

	// TODO return something other than form?
	return form, nil
}
