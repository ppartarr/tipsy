package correctors

import (
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Same is the identity function
func Same(password string) string {
	return password
}

// SwitchCaseFirstLetter switches the case of the first letter in the string to upper case
// swc-first
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
// swc-all
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
// rm-last
func RemoveLastChar(password string) string {
	lastCharRune, size := utf8.DecodeLastRuneInString(password)
	if lastCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the last rune")
		size = 0
	}
	return password[:len(password)-size]
}

// RemoveFirstChar removes the first character from the string
// rm-first
func RemoveFirstChar(password string) string {
	firstCharRune, size := utf8.DecodeRuneInString(password)
	if firstCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the first rune")
		size = 0
	}
	return password[size:]
}

// CapitalToUpper returns the password with every letter capitalised. For when users press shift-key instead of caps-lock
// cap2up
func CapitalToUpper(password string) string {
	if strings.Title(password) == password {
		return strings.ToUpper(password)
	}
	return password
}

// UpperToCapital returns the password with every letter capitalised. For when users press caps-lock instead of shift
// up2cap
func UpperToCapital(password string) string {
	if strings.ToUpper(password) == password {
		return strings.Title(strings.ToLower(password))
	}
	return password
}

// ConvertLastNumberToSymbol converts the last number to a symbol - TODO should depend on keyboard layout
// n2s-last
func ConvertLastNumberToSymbol(password string) string {
	lastCharRune, size := utf8.DecodeLastRuneInString(password)
	if lastCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the last rune")
		size = 0
	}
	if unicode.IsDigit(lastCharRune) {
		return password[:len(password)-size] + shiftSwitchMap[string(lastCharRune)]
	}
	return password
}

// SwitchShiftLastCharacter changes the last character according to the appropriate shift modifier
// sws-last1
func SwitchShiftLastCharacter(password string) string {
	lastCharRune, size := utf8.DecodeLastRuneInString(password)
	if lastCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the last rune")
		size = 0
	}
	if unicode.IsDigit(lastCharRune) || unicode.IsSymbol(lastCharRune) || unicode.IsPunct(lastCharRune) {
		return password[:len(password)-size] + shiftSwitchMap[string(lastCharRune)]
	} else if unicode.IsLetter(lastCharRune) {
		if unicode.IsUpper(lastCharRune) {
			return password[:len(password)-size] + string(unicode.ToLower(lastCharRune))
		} else if unicode.IsLower(lastCharRune) {
			return password[:len(password)-size] + string(unicode.ToUpper(lastCharRune))
		}
	}
	return password
}

// SwitchShiftLastNCharacters changes the last n characters according to the appropriate shift modifier
// sws-lastn
func SwitchShiftLastNCharacters(password string, n int) string {
	temp := password
	for i := 1; i <= n; i++ {
		// TODO improve this, it's very hacky...
		temp = SwitchShiftLastCharacter(temp[:len(temp)-(n-i)]) + password[len(password)-(n-i):]
	}

	return temp
}

// InverseRemoveLast appends every rune to the password
func InverseRemoveLast(password string) []string {
	edits := make([]string, 1)

	for _, letter := range letterRunes {
		// add every rune in every index
		edits = append(edits, password+string(letter))
	}

	return edits
}

// InverseRemoveFirst prepends every rune to the password
func InverseRemoveFirst(password string) []string {
	edits := make([]string, 1)

	for _, letter := range letterRunes {
		// add every rune in every index
		edits = append(edits, string(letter)+password)
	}

	return edits
}

// ConvertLastNumberToSymbol converts the last number to a symbol - TODO should depend on keyboard layout
// s2n-last
func ConvertLastSymbolToNumber(password string) string {
	lastCharRune, size := utf8.DecodeLastRuneInString(password)
	if lastCharRune == utf8.RuneError && (size == 0 || size == 1) {
		log.Fatal("Unable to decode the size of the last rune")
		size = 0
	}
	if unicode.IsSymbol(lastCharRune) {
		return password[:len(password)-size] + shiftSwitchMap[string(lastCharRune)]
	}
	return password
}

// AppendOne adds a 1 to the password
// add1_last
func AppendOne(password string) string {
	return password + "1"
}

// Keyboard layout dependent
var shiftSwitchMap = map[string]string{
	"`":  "~",
	"1":  "!",
	"2":  "@",
	"3":  "#",
	"4":  "$",
	"5":  "%",
	"6":  "^",
	"7":  "&",
	"8":  "*",
	"9":  "(",
	"0":  ")",
	"-":  "_",
	"=":  "+",
	"[":  "{",
	"]":  "}",
	"\\": "|",
	";":  ":",
	"'":  "\"",
	",":  "<",
	".":  ">",
	"/":  "?",
	"~":  "`",
	"!":  "1",
	"@":  "2",
	"#":  "3",
	"$":  "4",
	"%":  "5",
	"^":  "6",
	"&":  "7",
	"*":  "8",
	"(":  "9",
	")":  "0",
	"_":  "-",
	"+":  "=",
	"{":  "[",
	"}":  "]",
	"|":  "\\",
	":":  ";",
	"\"": "'",
	"<":  ",",
	">":  ".",
	"?":  "/",
}
