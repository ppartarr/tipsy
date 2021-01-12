todo

* implement rate limiting
    1. user clicks on password reset
    2. generate 64-bit token and store it alongside TTL & creation date
    3. email has been sent to account message (even if account doesn't exist to prevent user enumeration)
    4. user receives email with domain.com/reset/<token>
    5. backend validates token (user email matches token, TTL still valid)
    6. user fills in new password
    7. backend invalidates token and replaces old passwords hash with new password hash
* add zxcvbn by dropbox & only allow strong passwords to be registered
    * disable submit button until 3/4 zxcvbn is reached
* add config
    * make correctors modifiable
* implement personalised typo-correcion (read paper on Tuesday)
* if you are applying 3 checkers and rate-limiting at 10, try getting the 30th q most probable pasword in Blacklist to make sure there's never a decrease in security!
* send mail to quentin asking for ratio of success / failure of authentications without revealing the total number of attempts

how many passwords are pasted vs typed?

try get login failure rate from big companies (hint: Orange)
How many passwords are typed vs pasted?

## Links
seclists https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases