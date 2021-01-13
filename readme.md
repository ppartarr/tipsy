# todo

* don't allow users to register with password from blacklist
* add zxcvbn by dropbox & only allow strong passwords to be registered
    * disable submit button until 3/4 zxcvbn is reached
* add config
    * make correctors modifiable
* implement personalised typo-correcion (read paper on Tuesday)
* if you are applying 3 checkers and rate-limiting at 10, try getting the 30th q most probable pasword in Blacklist to make sure there's never a decrease in security!
* add all other correctors
* add error message on failed login
* write JavaScript to know if password is typed or pasted
* need to finish setting up HTML redirections & JavaScript
* check if user already has requested token before creating new one in PasswordReset before storeToken

## Important

* Send romke a project recap email and prepare question for Monday!
* Put authentication system online at typo.partarrieu.me

## How many passwords are pasted vs typed?

* test what banks allow password pasting
* PAM authentication


## Links
seclists https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases
go pake implementation: https://github.com/schollz/pake