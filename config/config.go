package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Server models the petze.yml main config file
type Server struct {
	SMTP *struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		From     string `yaml:"from"`
	} `yaml:"smtp"`

	Checker *struct {
		Always    bool              `yaml:"always"`
		Blacklist *BlacklistChecker `yaml:"blacklist"`
		Optimal   *OptimalChecker   `yaml:"optimal"`
	} `yaml:"checker"`

	Correctors []string `yaml:"correctors"`
	RateLimit  int      `yaml:"rate-limit"`
}

// BlacklistChecker represents blacklist checker
type BlacklistChecker struct {
	File string `yaml:"file"`
}

// OptimalChecker represents the optimal checker
type OptimalChecker struct {
	File                    string `yaml:"file"`
	QthMostProbablePassword int    `yaml:"qthMostProbablePassword"`
	Typos                   *Typos `yaml:"typos"`
}

// Typos represents the distribution of typos
type Typos struct {
	Same            int `yaml:"same"`
	Other           int `yaml:"other"`
	SwcAll          int `yaml:"swc-all"`
	Kclose          int `yaml:"kclose"`
	KeypressEdit    int `yaml:"keypress-edit"`
	RemoveLast      int `yaml:"rm-last"`
	SwitchFirst     int `yaml:"swc-first"`
	RemoveFirstChar int `yaml:"rm-firstc"`
	SwitchLast      int `yaml:"sws-last1"`
	Tcerror         int `yaml:"tcerror"`
	SwitchLastN     int `yaml:"sws-lastn"`
	Upncap          int `yaml:"upncap"`
	N2sLast         int `yaml:"n2s-last"`
	Cap2up          int `yaml:"cap2up"`
	Add1Last        int `yaml:"add1-last"`
}

// LoadServer loads the server configuration yaml
func LoadServer(configFile string) (server *Server, err error) {
	server = &Server{}

	// fill server struct from file
	err = load(configFile, &server)
	if err != nil {
		return server, load(configFile, &server)
	}

	// validate values
	err = server.IsValid()

	if err != nil {
		return server, err
	}

	return server, nil
}

// TODO write function to validate correctors

// Load a config interface from a file
func load(configFile string, target interface{}) error {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	yamlErr := yaml.UnmarshalStrict(configBytes, target)
	if yamlErr != nil {
		return errors.New("could not unmarshal yaml file " + configFile + " : " + yamlErr.Error())
	}
	return nil
}

// IsValid checks that the server config struct is valid
func (s *Server) IsValid() error {
	// check that only a single corrector is defined in the config
	err := s.validateNumberOfCheckers()
	if err != nil {
		return err
	}

	// TODO add validation

	return nil
}

func (s *Server) validateNumberOfCheckers() error {
	numberOfCheckers := 0

	if s.Checker.Always {
		numberOfCheckers++
	}

	if s.Checker.Blacklist != nil {
		numberOfCheckers++
	}

	if s.Checker.Optimal != nil {
		numberOfCheckers++
	}

	if numberOfCheckers != 1 {
		return errors.New("more than one checker is defined in config - use only one of: always, blacklist, optimal")
	}

	return nil
}
