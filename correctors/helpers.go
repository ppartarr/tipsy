package correctors

import (
	"log"
	"sort"
)

// Corrector constants
const (
	SwitchAll          = "swc-all"
	RemoveLast         = "rm-last"
	SwitchFirst        = "swc-first"
	RemoveFirst        = "rm-first"
	SwitchLast         = "sws-last1"
	SwitchLastN        = "sws-lastn"
	UpperNCapital      = "upncap"
	NumberToSymbolLast = "n2s-last"
	Capital2Upper      = "cap2up"
	AddOneLast         = "add1-last"
)

var allCorrectors = []string{
	SwitchAll,
	RemoveLast,
	SwitchFirst,
	RemoveFirst,
	SwitchLast,
	SwitchLastN,
	UpperNCapital,
	NumberToSymbolLast,
	Capital2Upper,
	AddOneLast,
}

// ApplyCorrectionFunction applies the appropriate corrector function given it's config name
func ApplyCorrectionFunction(corrector string, password string) string {
	// log.Println(corrector)
	// log.Println(password)

	switch corrector {
	case SwitchAll:
		return SwitchCaseAll(password)
	case RemoveLast:
		return RemoveLastChar(password)
	case SwitchFirst:
		return SwitchCaseFirstLetter(password)
	case RemoveFirst:
		return RemoveFirstChar(password)
	case SwitchLast:
		return SwitchShiftLastCharacter(password)
	case SwitchLastN:
		return SwitchShiftLastNCharacters(password, 2)
	case UpperNCapital:
		return UpperToCapital(password)
	case NumberToSymbolLast:
		return ConvertLastNumberToSymbol(password)
	case Capital2Upper:
		return CapitalToUpper(password)
	case AddOneLast:
		return AppendOne(password)
	}

	log.Fatal("corrector unknown:", corrector)
	return password
}

// ApplyInverseCorrectionFunction applies the appropriate corrector function given it's config name
func ApplyInverseCorrectionFunction(corrector string, password string) []string {
	inverse := make([]string, 1)

	switch corrector {
	case SwitchAll:
		inverse = append(inverse, SwitchCaseAll(password))
	case RemoveLast:
		edits := InverseRemoveLast(password)
		for _, edit := range edits {
			inverse = append(inverse, edit)
		}
	case SwitchFirst:
		inverse = append(inverse, SwitchCaseFirstLetter(password))
	case RemoveFirst:
		edits := InverseRemoveFirst(password)
		for _, edit := range edits {
			inverse = append(inverse, edit)
		}
	case SwitchLast:
		inverse = append(inverse, SwitchShiftLastCharacter(password))
	case SwitchLastN:
		inverse = append(inverse, SwitchShiftLastNCharacters(password, 2))
	case UpperNCapital:
		inverse = append(inverse, UpperToCapital(password))
	case NumberToSymbolLast:
		inverse = append(inverse, ConvertLastSymbolToNumber(password))
	case Capital2Upper:
		inverse = append(inverse, CapitalToUpper(password))
	case AddOneLast:
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

// TODO extend to use any rune instead of US alphanumerics
var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?")
