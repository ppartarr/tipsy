package typtop

import (
	"crypto/rsa"

	"github.com/ppartarr/tipsy/config"
)

// CheckerService represents a typtop checker service...
type CheckerService struct {
	config        *config.TypTopChecker
	typoFrequency map[string]int
}

// NewCheckerService initialises the CheckerService
func NewCheckerService(typtopConfig *config.TypTopChecker, typoFrequency map[string]int) (checker *CheckerService) {
	return &CheckerService{
		config:        typtopConfig,
		typoFrequency: typoFrequency,
	}
}

// User represents a TypTop user
type User struct {
	ID            int             `json:"id"`
	Email         string          `json:"email"`
	LoginAttempts int             `json:"loginAttempts"`
	PrivateKey    *rsa.PrivateKey `json:"privateKey"`
	State         *State          `json:"typtopState"`
}
