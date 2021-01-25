package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckBlacklist(t *testing.T) {
	checker := NewChecker(testTypos, topCorrectors)
	ball := checker.CheckBlacklist("password!", []string{"password"})
	if !assert.ElementsMatch(t, ball, []string{"Password!", "PASSWORD!"}) {
		t.Error("ball should be the set of strings containing the output of the correctors, unless it's in the blacklist")
	}
}
