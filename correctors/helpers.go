package correctors

import (
	"log"
	"sort"
)

// ApplyCorrectionFunction applies the appropriate corrector function given it's config name
func ApplyCorrectionFunction(corrector string, password string) string {
	// log.Println(corrector)
	// log.Println(password)

	switch corrector {
	case "swc-all":
		return SwitchCaseAll(password)
	case "rm-last":
		return RemoveLastChar(password)
	case "swc-first":
		return SwitchCaseFirstLetter(password)
	case "rm-first":
		return RemoveFirstChar(password)
	case "sws-last1":
		return SwitchShiftLastCharacter(password)
	case "sws-lastn":
		return SwitchShiftLastNCharacters(password, 2)
	case "upncap":
		return UpperToCapital(password)
	case "n2s-last":
		return ConvertLastNumberToSymbol(password)
	case "cap2up":
		return CapitalToUpper(password)
	case "add1-last":
		return AppendOne(password)
	}

	log.Fatal("corrector unknown:", corrector)
	return password
}

// ApplyInverseCorrectionFunction applies the appropriate corrector function given it's config name
func ApplyInverseCorrectionFunction(corrector string, password string) []string {
	inverse := make([]string, 1)

	switch corrector {
	case "swc-all":
		inverse = append(inverse, SwitchCaseAll(password))
	case "rm-last":
		edits := InverseRemoveLast(password)
		for _, edit := range edits {
			inverse = append(inverse, edit)
		}
	case "swc-first":
		inverse = append(inverse, SwitchCaseFirstLetter(password))
	case "rm-first":
		edits := InverseRemoveFirst(password)
		for _, edit := range edits {
			inverse = append(inverse, edit)
		}
	case "sws-last1":
		inverse = append(inverse, SwitchShiftLastCharacter(password))
	case "sws-lastn":
		inverse = append(inverse, SwitchShiftLastNCharacters(password, 2))
	case "upncap":
		inverse = append(inverse, CapitalToUpper(password))
	case "n2s-last":
		inverse = append(inverse, ConvertLastSymbolToNumber(password))
	case "cap2up":
		inverse = append(inverse, UpperToCapital(password))
	case "add1-last":
		inverse = append(inverse, RemoveLastChar(password))
	}

	return inverse
}

// KeyValue reporesents a map as a slice
type KeyValue struct {
	Key   string
	Value int
}

// GetNBestCorrectors returns the n best correctors in order, determined by the typo frequency
func GetNBestCorrectors(n int, typoFrequency map[string]int) []string {

	// there are only 10 correctors implemented
	if n > 10 {
		log.Fatal("there are only 10 correctors")
	}

	nBestCorrectors := make([]string, 0)

	ss := ConvertMapToSortedSlice(typoFrequency)

	// add corrector to slice
	for i := 0; i < len(typoFrequency); i++ {
		// some typos aren't correction functions e.g. keypress-edit, tcerror, other
		if StringInSlice(ss[i].Key, allCorrectors) && len(nBestCorrectors) < n {
			nBestCorrectors = append(nBestCorrectors, ss[i].Key)
		}
	}

	nBestCorrectors = DeleteEmpty(nBestCorrectors)
	return nBestCorrectors
}

// ConvertMapToSortedSlice converts a map to a slice and sorts it by value
func ConvertMapToSortedSlice(in map[string]int) []KeyValue {
	var ss []KeyValue

	for k, v := range in {
		ss = append(ss, KeyValue{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

// GetBall returns the passwords in the ball given a slice of correctors
func GetBall(password string, correctors []string) []string {
	ball := make([]string, len(correctors))
	for index, corrector := range correctors {
		correctedPassword := ApplyCorrectionFunction(corrector, password)
		if correctedPassword != password {
			ball[index] = correctedPassword
		}
	}
	ball = DeleteEmpty(ball)
	return ball
}

// GetBallWithCorrectionType returns the ball with the correction type string
func GetBallWithCorrectionType(password string, correctors []string) map[string]string {
	var ballWithCorrectorName = make(map[string]string)

	for _, corrector := range correctors {
		correctedPassword := ApplyCorrectionFunction(corrector, password)
		if correctedPassword != password {
			ballWithCorrectorName[correctedPassword] = corrector
		}
	}

	return ballWithCorrectorName
}

// DeleteEmpty remove all empty strings from a slice
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// StringInSlice returns true if s in is list
func StringInSlice(s string, list []string) bool {
	success := false
	for _, value := range list {
		if value == s {
			success = true
		}
	}
	return success
}

var allCorrectors = []string{
	"swc-all",
	"rm-last",
	"swc-first",
	"rm-first",
	"sws-last1",
	"sws-lastn",
	"upncap",
	"n2s-last",
	"cap2up",
	"add1-last",
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?")
