package typtop

import (
	"crypto/rsa"

	"github.com/ppartarr/tipsy/config"
)

// Checker represents a typtop checker service...
type Checker struct {
	config        *config.TypTopChecker
	typoFrequency map[string]int
}

// NewChecker initialises the Checker
func NewChecker(typtopConfig *config.TypTopChecker, typoFrequency map[string]int) (checker *Checker) {
	return &Checker{
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
