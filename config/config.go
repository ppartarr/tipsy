package config

import "log"

const (
	configFile = "tipsy.yml"
)

type Checker struct (
	AlwaysChecker *AlwaysChecker
	BlacklistChecker *BlacklistChecker
	OptimalChecker *OptimalChecker
)

type AlwaysChecker (
)

type BlacklistChecker (
	File string `yaml:"file`
)

type OptimalChecker struct (
	File string `yaml:"file`
)

type Correctors struct (
	Correctors []string `yaml:"correctors`
)

func Correctors() []string {
	return correctors
}

func SetCorrectors(correctors []string) {
	for _, corrector := range correctors {
		switch corrector {
			case 
			"same",
			"swc-all",
			"swc-first",
			"rm-last",
			"rm-first": return true
		}
		log.Fatal("correctors must be a list of strings with accepted values from: swc-all, swc-first, rm-last, rm-first")
	}
}