package checkers

// Checker represents a checker service...
type Checker struct {
	TypoFrequency map[string]int
	Correctors    []string
}

// NewChecker initialises the Checker
func NewChecker(typoFrequency map[string]int, correctors []string) (checker *Checker) {
	return &Checker{
		TypoFrequency: typoFrequency,
		Correctors:    correctors,
	}
}
