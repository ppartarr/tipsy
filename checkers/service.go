package checkers

// CheckerService represents a checker service...
type CheckerService struct {
	TypoFrequency map[string]int
	Correctors    []string
}

// NewCheckerService initialises the CheckerService
func NewCheckerService(typoFrequency map[string]int, correctors []string) (checker *CheckerService) {
	return &CheckerService{
		TypoFrequency: typoFrequency,
		Correctors:    correctors,
	}
}
