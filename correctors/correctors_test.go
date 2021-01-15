package correctors

import "testing"

func TestSame(t *testing.T) {
	if Same("password") != "password" {
		t.Errorf("should be the identity function for string")
	}
}

func TestSwitchCaseFirstLetter(t *testing.T) {
	if SwitchCaseFirstLetter("password") != "Password" {
		t.Errorf("should switch the case of the first letter")
	}

	if SwitchCaseFirstLetter("Password") != "password" {
		t.Errorf("should switch the case of the first letter")
	}

	if SwitchCaseFirstLetter("4assword") != "4assword" {
		t.Errorf("should ONLY switch the case of the first letter (not symbols)")
	}

	if SwitchCaseFirstLetter("") != "" {
		t.Errorf("should handle empty strings")
	}
}

func TestSwitchCaseAll(t *testing.T) {
	if SwitchCaseAll("password") != "PASSWORD" {
		t.Errorf("should switch the case of all letters in the string")
	}

	if SwitchCaseAll("PASSWORD") != "password" {
		t.Errorf("should switch the case of all letters in the string")
	}

	if SwitchCaseAll("password1") != "PASSWORD1" {
		t.Errorf("should switch the case of all letters in the string")
	}
}

func TestRemoveLastChar(t *testing.T) {
	if RemoveLastChar("password") != "passwor" {
		t.Errorf("should remove the last character in the string")
	}
	// test invalid utf-8
	// if RemoveLastChar("") != "" {
	// 	t.Errorf("should return the input string if it fails to decode the rune")
	// }
}

func TestRemoveFirstChar(t *testing.T) {
	if RemoveFirstChar("password") != "assword" {
		t.Errorf("RemoveLastChat should remove the first character in the string")
	}
	// test invalid utf-8
	// if RemoveFirstChar("") != "" {
	// 	t.Errorf("should return the input string if it fails to decode the rune")
	// }
}

func TestCapitalToUpper(t *testing.T) {
	if CapitalToUpper("password") != "password" {
		t.Errorf("should capitalise every letter ONLY if the password begins with a capital letter")
	}

	if CapitalToUpper("Password") != "PASSWORD" {
		t.Errorf("should capitalise every letter if the password starts with a capital")
	}

	if CapitalToUpper("PASSWORD") != "PASSWORD" {
		t.Errorf("should capitalise every letter ONLY if the password begins with a capital letter")
	}
}

func TestUpperToCapital(t *testing.T) {
	if UpperToCapital("PASSWORD") != "Password" {
		t.Errorf("should switch the capitalisation except for the first char")
	}

	if UpperToCapital("Password") != "Password" {
		t.Errorf("should ONLY switch the capitalisation if the entire string is capitalised")
	}

	if UpperToCapital("password") != "password" {
		t.Errorf("should ONLY switch the capitalisation if the entire string is capitalised")
	}
}

func TestConvertLastNumberToSymbol(t *testing.T) {
	if ConvertLastNumberToSymbol("password1") != "password!" {
		t.Errorf("should convert the last number to a symbol")
	}

	if ConvertLastNumberToSymbol("password!") != "password!" {
		t.Errorf("should ONLY convert the last NUMBER to a symbol")
	}
}

func TestSwitchShiftLastCharacter(t *testing.T) {
	if SwitchShiftLastCharacter("password") != "passworD" {
		t.Errorf("should capitalise the last character if it's a letter")
	}

	if SwitchShiftLastCharacter("passworD") != "password" {
		t.Errorf("should capitalise the last character if it's a letter")
	}

	if SwitchShiftLastCharacter("password1") != "password!" {
		t.Errorf("should replace the last number with the appropriate symbol (determined by the shift modifier)")
	}

	if SwitchShiftLastCharacter("password!") != "password1" {
		t.Errorf("should replace the last symbol with the appropriate number (determined by the shift modifier)")
	}
}

func TestSwitchShiftLastNCharacters(t *testing.T) {
	if SwitchShiftLastNCharacters("password", 3) != "passwORD" {
		t.Errorf("should capitalise the last character if it's a letter")
	}

	if SwitchShiftLastNCharacters("passworD", 3) != "passwORd" {
		t.Errorf("should capitalise the last character if it's a letter")
	}

	if SwitchShiftLastNCharacters("password123", 3) != "password!@#" {
		t.Errorf("should replace the last number with the appropriate symbol (determined by the shift modifier)")
	}

	if SwitchShiftLastNCharacters("passwoRd!", 3) != "passworD1" {
		t.Errorf("should replace the last symbol with the appropriate number (determined by the shift modifier)")
	}
}

func TestAppendOne(t *testing.T) {
	if AppendOne("password") != "password1" {
		t.Errorf("should add 1 to the end of the string")
	}
}
