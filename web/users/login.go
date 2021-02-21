package users

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mcnijman/go-emailaddress"
	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/checkers/typtop"
	"github.com/ppartarr/tipsy/web/session"
)

// LoginForm represents a login form
type LoginForm struct {
	Email    string
	Password string
	Pasted   string
	NoJS     string
	Errors   map[string]string
}

// Validate checks that the fields in the login form are set
func (form *LoginForm) Validate() bool {
	form.Errors = make(map[string]string)

	if strings.TrimSpace(form.Password) == "" {
		form.Errors["Login"] = "Username and password incorrect"
	}

	// validate email
	log.Println("validating email address")
	_, err := emailaddress.Parse(form.Email)
	if err != nil {
		form.Errors["Login"] = "Username and password incorrect"
	}

	// log pasted value
	log.Println("user:", form.Email, " password:", form.Password, " pasted:", form.Pasted, " nojs:", form.NoJS)

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
		Pasted:   r.PostFormValue("pasted"),
		NoJS:     r.PostFormValue("nojs"),
	}
	if form.Validate() == false {
		log.Println(form.Errors)
		return form, errors.New("you must submit a valid form")
	}

	if userService.config.Checker.TypTop == nil {

		// get user from email
		user, err := userService.getUser(form.Email)
		if err != nil {
			log.Println("could not get user " + form.Email + ": " + err.Error())
			form.Errors["Login"] = "Username and password incorrect"
			return form, errors.New("you must submit a valid form")
		}

		// use one of always, blacklist, optimal
		if userService.config.Checker.Always {
			log.Println("using always checker")
			// check password
			ball := userService.checker.CheckAlways(form.Password)
			if !checkPasswordAndBall(form.Password, ball, user.PasswordHash) {

				// increment login attempts
				err = userService.checkLoginAttempts(w, r, user)
				if err != nil {
					return form, errors.New("failed to check login attemps")
				}
			}
		} else if userService.config.Checker.Blacklist != nil {
			log.Println("using blacklist checker")
			blacklist := checkers.LoadBlacklist(userService.config.Checker.Blacklist.File)
			ball := userService.checker.CheckBlacklist(form.Password, blacklist)

			// if password check fails, increment login attempts
			if !checkPasswordAndBall(form.Password, ball, user.PasswordHash) {

				// increment login attempts
				err = userService.checkLoginAttempts(w, r, user)
				if err != nil {
					return form, errors.New("failed to check login attemps")
				}
			}
		} else if userService.config.Checker.Optimal != nil {
			log.Println("using optimal checker")
			frequencyBlacklist := checkers.LoadFrequencyBlacklist(userService.config.Checker.Optimal.File, userService.config.MinPasswordLength)
			ball := userService.checker.CheckOptimal(form.Password, frequencyBlacklist, userService.config.Checker.Optimal.QthMostProbablePassword)

			// if password check fails, increment login attempts
			if !checkPasswordAndBall(form.Password, ball, user.PasswordHash) {
				// increment login attempts
				err = userService.checkLoginAttempts(w, r, user)
				if err != nil {
					return form, errors.New("failed to check login attemps")
				}
			}
		}

		log.Println("successfully logged in")

		// init session
		_, err = session.SetUserID(w, r, strconv.Itoa(user.ID))
		if err != nil {
			return form, errors.New("could not create a session for user " + form.Email + ": " + err.Error())
		}

	} else {
		log.Println("using typtop checker")

		typtopUser, err := userService.getTypTopUser(form.Email)
		if err != nil {
			log.Println("could not get typtop user " + form.Email + ": " + err.Error())
			form.Errors["Login"] = "Username and password incorrect"
			return form, errors.New("you must submit a valid form")
		}

		success, typtopState := userService.typtop.Login(typtopUser.State, form.Password, typtopUser.PrivateKey)

		if !success {
			// increment login attempts
			err = userService.checkTypTopLoginAttempts(w, r, typtopUser)
			if err != nil {
				return form, errors.New("failed to check login attemps")
			}
		}

		// update user state
		typtopUser.State = typtopState

		userService.updateTypTopUser(typtopUser)

		log.Println("successfully logged in")

		// init session
		_, err = session.SetUserID(w, r, strconv.Itoa(typtopUser.ID))
		if err != nil {
			return form, errors.New("could not create a session for typtop user " + form.Email + ": " + err.Error())
		}
	}

	return form, nil
}

// Check the submitted password first, then remainder of the ball
// this way timing attacks only tell if a corrector is used
func checkPasswordAndBall(submittedPassword string, ball []string, registeredPasswordHash string) bool {
	if CheckPasswordHash(submittedPassword, registeredPasswordHash) {
		return true
	}

	success := false
	for _, password := range ball {
		if CheckPasswordHash(password, registeredPasswordHash) {
			success = true
		}
	}

	return success
}

func (userService *UserService) checkLoginAttempts(w http.ResponseWriter, r *http.Request, user *User) error {
	// increment login attempts
	if user.LoginAttempts < userService.config.Web.Login.RateLimit {
		userService.incrementLoginAttempts(user)
	} else {
		log.Println("user has attempted to login too many times")
		go func() {
			userService.PasswordRecovery(w, r)
		}()
		return errors.New("the limit of login attempts has been reached, please reset your password via the mail provided in the link " + r.Form["email"][0])
	}

	return errors.New("wrong password for user " + r.Form["email"][0])
}

func (userService *UserService) checkTypTopLoginAttempts(w http.ResponseWriter, r *http.Request, user *typtop.User) error {
	// increment login attempts
	if user.LoginAttempts < userService.config.Web.Login.RateLimit {
		userService.incrementTypTopLoginAttempts(user)
	} else {
		log.Println("user has attempted to login too many times")
		go func() {
			userService.PasswordRecovery(w, r)
		}()
		return errors.New("the limit of login attempts has been reached, please reset your password via the mail provided in the link " + r.Form["email"][0])
	}

	return errors.New("wrong password for user " + r.Form["email"][0])
}
