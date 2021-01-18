# todo

* don't allow users to register with password from blacklist
* add zxcvbn by dropbox & only allow strong passwords to be registered
    * disable submit button until 3/4 zxcvbn is reached
* implement personalised typo-correcion (read paper on Tuesday)
* if you are applying 3 checkers and rate-limiting at 10, try getting the 30th q most probable pasword in Blacklist to check if there's a decrease in security compared to the exact checker !
* add error message on failed login
* need to finish setting up HTML redirections & JavaScript
* check if user already has requested token before creating new one in PasswordReset before storeToken
* ask Phil about encoding before saving to db - tutorial said json but that feels dodgy

## Important

* Put authentication system online at typo.partarrieu.me

## Future work

* make correctors work with different keyboard layouts

## How many passwords are pasted vs typed?

* test what banks allow password pasting
* PAM authentication


## Links
seclists https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases
go pake implementation: https://github.com/schollz/pake