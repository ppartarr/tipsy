package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTypos = map[string]int{
	"same":          90234,
	"other":         1918,
	"swc-all":       1698,
	"kclose":        1385,
	"keypress-edit": 1000,
	"rm-last":       382,
	"swc-first":     209,
	"rm-first":      55,
	"sws-last":      19,
	"tcerror":       18,
	"sws-lastn":     14,
	"upncap":        13,
	"n2s-last":      9,
	"cap2up":        5,
	"add1-last":     5,
}

var topCorrectors = []string{"swc-all", "rm-last", "swc-first"}

func TestCheckAlways(t *testing.T) {
	checker := NewChecker(testTypos, topCorrectors)
	ball := checker.CheckAlways("test")
	if !assert.ElementsMatch(t, ball, []string{"Test", "TEST", "tes"}) {
		t.Error("ball should be the set of strings containing the output of the correctors applied to the password")
	}
}
