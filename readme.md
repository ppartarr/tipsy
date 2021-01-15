# todo

* personalised typo correction
    * implement encrypted wait list: W
        * public-key pk encryption of recent incorrect password submissions that are not the registered password or one of the typos in the typo-cache e.g. for registered password p and incorrect submissions w0, w1, ..., wn such that w != w0 and for typo cache T, w ∉ T
            * { {w0}pk, {w1}pk, ..., {wn}pk }
                * ensure that p != w, for any w in {w0, w1, ..., wn}
                * given a typo cache T, w ∉ T, for any {w0, w1, ..., wn}
                * has a configurable length that should match the rateLimit in practice
            * uses least recently used (LRU) eviction policy
        * private-key sk is encrypted using the registered password:
            * C = E(sk, p)
            * p is the user's registered password
            * sk is the private-key (secret key)
    * implement typo cache: T
        * set of strings that the user is allowed to authenticate with (registed passwords + accepted typos of password) e.g. for password w and accepted submissions wa, wb, ..., wn
            * { E(sk, wa), E(sk, wb), ..., E(sk, wn) }
            * has configurable length
        * cache initialisation
        * cache update policy (configurable policy e.g. LFU, PLFU etc)
        * after every change in typo cache, randomly permute order of cached typos
    * write isTypo() based on the Damerau-Levenshtein distance
    * state = (publicKey, cachingSchemeState, typoCache, waitList, nextWaitListEntry)
    * admissable typos must respect 3 (configurable) conditions:
        * Damerau-Levenshtein < 1
        * strength estimation of typo password >= 10   (ensures easily guessed passwords are never cached using zxcvbn)
        * strength estimation of typo password >= strength estimation of password - 3  (prevent caching of typos significantly more guessable than the real password)

* don't allow users to register with password from blacklist
* add zxcvbn by dropbox & only allow strong passwords to be registered
    * disable submit button until 3/4 zxcvbn is reached
* finish writing all correctors
* implement personalised typo-correcion (read paper on Tuesday)
* if you are applying 3 checkers and rate-limiting at 10, try getting the 30th q most probable pasword in Blacklist to check if there's a decrease in security compared to the exact checker !
* add all other correctors
* add error message on failed login
* write JavaScript to know if password is typed or pasted
* need to finish setting up HTML redirections & JavaScript
* check if user already has requested token before creating new one in PasswordReset before storeToken

## Important

* Send romke a project recap email and prepare question for Monday!
* Put authentication system online at typo.partarrieu.me
* answer ethical team question email!

## How many passwords are pasted vs typed?

* test what banks allow password pasting
* PAM authentication


## Links
seclists https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases
go pake implementation: https://github.com/schollz/pake