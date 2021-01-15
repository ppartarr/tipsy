package checkers

// CheckerService represents a checker service...
type CheckerService struct {
	TypoFrequency      map[string]int
	NumberOfCorrectors int
}

// NewCheckerService initialises the CheckerService
func NewCheckerService(typoFrequency map[string]int, numberOfCorrectors int) (checker *CheckerService) {
	return &CheckerService{
		TypoFrequency:      typoFrequency,
		NumberOfCorrectors: numberOfCorrectors,
	}
}
