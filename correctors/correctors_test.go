package correctors

import "testing"

func TestSwitchCaseFirstLetter(t *testing.T) {
	if SwitchCaseFirstLetter("password") != "Password" {
		t.Errorf("SwitchCaseFirstLetters should switch the case of the first letter")
	}

	if SwitchCaseFirstLetter("Password") != "password" {
		t.Errorf("SwitchCaseFirstLetters should switch the case of the first letter")
	}

	if SwitchCaseFirstLetter("4assword") != "4assword" {
		t.Errorf("SwitchCaseFirstLetters should ONLY switch the case of the first letter (not symbols)")
	}

	if SwitchCaseFirstLetter("") != "" {
		t.Errorf("SwitchCaseFirstLetters should handle empty strings")
	}
}

func TestSwitchCaseAll(t *testing.T) {
	if SwitchCaseAll("password") != "PASSWORD" {
		t.Errorf("SwitchCaseAll should switch the case of all letters in the string")
	}

	if SwitchCaseAll("PASSWORD") != "password" {
		t.Errorf("SwitchCaseAll should switch the case of all letters in the string")
	}

	if SwitchCaseAll("password1") != "PASSWORD1" {
		t.Errorf("SwitchCaseAll should switch the case of all letters in the string")
	}
}

func TestRemoveLastChar(t *testing.T) {
	if RemoveLastChar("password") != "passwor" {
		t.Errorf("RemoveLastChar should remove the last character in the string")
	}
	// test invalid utf-8
	// if RemoveLastChar("") != "" {
	// 	t.Errorf("RemoveLastChat should return the input string if it fails to decode the rune")
	// }
}

func TestRemoveFirstChar(t *testing.T) {
	if RemoveFirstChar("password") != "assword" {
		t.Errorf("RemoveLastChat should remove the first character in the string")
	}
	// test invalid utf-8
	// if RemoveFirstChar("") != "" {
	// 	t.Errorf("RemoveLastChat should return the input string if it fails to decode the rune")
	// }
}
