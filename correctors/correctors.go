package correctors

import (
	"log"
	"unicode"
	"unicode/utf8"
)

// Same is the identity function
func Same(password string) string {
	return password
}

// SwitchCaseFirstLetter switches the case of the first letter in the string to upper case
func SwitchCaseFirstLetter(password string) string {
	// we avoid toTitle in case the password contains a space
	for index, value := range password {
		if unicode.IsLower(value) {
			return string(unicode.ToUpper(value)) + password[index+1:]
		}
		return string(unicode.ToLower(value)) + password[index+1:]
	}

	// TODO should probably raise error here is password length is 0
	return ""
}

// SwitchCaseAll switches the case of all the letters in the string
func SwitchCaseAll(password string) string {
	newPassword := ""
	for _, value := range password {
		if unicode.IsLower(value) {
			newPassword = newPassword + string(unicode.ToUpper(value))
		} else if unicode.IsUpper(value) {
			newPassword = newPassword + string(unicode.ToLower(value))
		} else {
			newPassword = newPassword + string(value)
		}
	}
	return newPassword
}

// RemoveLastChar removes the last character from the string
func RemoveLastChar(password string) string {
	lastCharRune, size := utf8.DecodeLastRuneInString(password)
	if lastCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the last rune")
		size = 0
	}
	return password[:len(password)-size]
}

// RemoveFirstChar removes the first character from the string
func RemoveFirstChar(password string) string {
	firstCharRune, size := utf8.DecodeRuneInString(password)
	if firstCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the first rune")
		size = 0
	}
	return password[size:]
}
