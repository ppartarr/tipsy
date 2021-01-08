package checkers

// CheckBlacklist uses a blacklist of high-probability passwords. It checks the password or any password in the ball
func CheckBlacklist(password string, blacklist []string) bool {
	// TODO get the password from db
	registeredPassword := "password"

	var ball []string = getBall(password)

	// check the submitted password first
	if password == registeredPassword {
		return true
	}

	for _, value := range ball {
		// check password in the ball only if it isn't in the blacklist
		if !stringInSlice(value, blacklist) {
			return registeredPassword == value
		}
	}

	return false
}
