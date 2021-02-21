# Pasted.js

Collect logs to measure the utility of correcting typos on your service! This will allow us to measure how many passwords are pasted (we also check how many clients have JS disabled to refine our calculation).


|          | Bitwarden | LastPass | 1Password |
| ---      | ---       | ---      | ---       |
| Chrome   | ✅        |          |           |
| Firefox  |           |          |           |
| Opera    |           |          |           |
| Safari   |           |          |           |
| Brave    |           |          |           |


Submit an issue if you'd like add a browser or password manager to the list!


### Client side
We will need to modify the login form to:
1. add a hidden checkbox
2. modify the password input field to set the checkbox with `isPasted()`
3. include pasted.js (the script can also be included inline)

```html
<!-- example login.html -->
<!doctype html>
<html>
    <body>
        <form action="/login" method="post" autocomplete="off">
            <input id="email" type="text" placeholder="Email" name="email" required>
            <input id="password" type="password" placeholder="Password" name="password" required onpaste="isPasted()">
            <input id="pasted" type="checkbox" name="pasted" hidden checked=false value="true">
            <button type="submit"><span>Login</span></button>
        </form>
    </body>
    <footer>
        <script type="text/javascript" src="/static/scripts/pasted.js"></script>
    </footer>
</html>
```

Not all browsers and password managers work in the same way so there isn't one reliable event that we can listen to. Instead we want to support as many browsers and password managers as possible by:
1. detecting paste events in the password field
2. detecting autofill events in the password field

```js
// pasted.js (up to date version is in static/scripts/)
"use strict";

isAutoFilledListener();

function isPasted() {
    const pastedCheckbox = document.getElementById("pasted");
    pastedCheckbox.setAttribute("checked", true);
    console.log(pastedCheckbox.getAttribute("checked"));
}

function isAutoFilledListener() {
    // get your password input element
    const passwordElement = document.getElementById("password"); 

    // function to execute when input event happens
    function inputListenerFunc(event) { 
        if (
            event.target.value !== '' &&
            event.inputType === undefined &&
            event.data === undefined &&
            event.dataTransfer === undefined &&
            event.isComposing === undefined
        ) {
            // code to run when password manager auto-fill has occured
            isPasted();
        }
    }

    passwordElement.addEventListener('input', inputListenerFunc);
}
```
### Example server side
For every login request we simply want to log the value of the pasted checkbox. This will depend on the programming language, framework and logger that you use. The submitted client form will contain an email, password, and pasted field.

In go it will look like this:
```go
// LoginForm represents a login form
type LoginForm struct {
	Email    string
	Password string
	Pasted   string
	Errors   map[string]string
}

func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) error {
	// check that a valid form was submitted
	r.ParseForm()
	form = &LoginForm{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		Pasted:   r.PostFormValue("pasted"),
	}

    // log pasted value
	log.Println("user:", form.Email, " password:", form.Password, " pasted:", form.Pasted)

    // authenticate user
}
```

Here are some example logs:
```txt
# manually typed password
2021/02/21 11:38:46 user: test@partarrieu.me  password: password  pasted: 

# autofilled password
2021/02/21 11:36:46 user: test@partarrieu.me  password: password  pasted: true
```

DO NOT LOG PASSWORDS IN PRODUCTION! 🙏


## Adding JavaScript detection

To get more accurate results, simply add this to your login form:
### Client side
Simply 
```html
<noscript>
    <input type="checkbox" name="nojs" id="nojs" value="true" hidden checked>
</noscript>
```

### Example server side
```go
// LoginForm represents a login form
type LoginForm struct {
	Email    string
	Password string
	Pasted   string
    NoJS     string
	Errors   map[string]string
}

func (userService *UserService) Login(w http.ResponseWriter, r *http.Request) error {
	// check that a valid form was submitted
	r.ParseForm()
	form = &LoginForm{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		Pasted:   r.PostFormValue("pasted"),
        NoJS:     r.PostFormValue("nojs"),
	}

    // log pasted value
	log.Println("user:", form.Email, " password:", form.Password, " pasted:", form.Pasted, " nojs:", form.NoJS)

    // authenticate user
}
```

## Links
* https://stackoverflow.com/questions/11708092/detecting-browser-autofill
* https://js.plainenglish.io/how-to-detect-passwords-managers-usage-on-a-website-using-javascript-97fc1dff5c4a