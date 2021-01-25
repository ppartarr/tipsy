package config

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// Server models the petze.yml main config file
type Server struct {
	Smtp *SMTP `yaml:"smtp"`

	Checker *Checker `yaml:"checker"`

	Typos      map[string]int `yaml:"typos"`
	Correctors []string       `yaml:"correctors"`
	Web        *Web           `yaml:"web"`
}

type Checker struct {
	Always    bool              `yaml:"always"`
	Blacklist *BlacklistChecker `yaml:"blacklist"`
	Optimal   *OptimalChecker   `yaml:"optimal"`
	TypTop    *TypTopChecker    `yaml:"typtop"`
}

type SMTP struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
}

// Web is the config for the webserver
type Web struct {
	Register *Register `yaml:"register"`
	Login    *Login    `yaml:"login"`
	Reset    *Reset    `yaml:"reset"`
}

// Register represents the registration page & API
type Register struct {
	Blacklist string `yaml:"blacklist"`
	Zxcvbn    int    `yaml:"zxcvbn"`
}

// Login represents the login page & API
type Login struct {
	RateLimit       int           `yaml:"rateLimit"`
	SessionValidity time.Duration `yaml:"sessionValidity"`
}

// Reset represents the reset & API
type Reset struct {
	TokenValidity time.Duration `yaml:"tokenValidity"`
}

// BlacklistChecker represents blacklist checker
type BlacklistChecker struct {
	File string `yaml:"file"`
}

// OptimalChecker represents the optimal checker
type OptimalChecker struct {
	File                    string `yaml:"file"`
	QthMostProbablePassword int    `yaml:"qthMostProbablePassword"`
}

// TypTopChecker represents the typtop checker
type TypTopChecker struct {
	PublicKeyEncryption     PublicKeyEncryption     `yaml:"pke"`
	PasswordBasedEncryption PasswordBasedEncryption `yaml:"pbe"`
	EditDistance            int                     `yaml:"editDistance"`
	Zxcvbn                  int                     `yaml:"zxcvbn"`
	TypoCache               TypoCache               `yaml:"typoCache"`
	WaitList                WaitList                `yaml:"waitList"`
}

// PublicKeyEncryption represents the PKE scheme used for typtop
type PublicKeyEncryption struct {
	KeyLength int `yaml:"keyLength"`
}

// PasswordBasedEncryption represents the PBE scheme used for typtop
type PasswordBasedEncryption struct {
	KeyLength int `yaml:"keyLength"`
}

// TypoCache represents the config for the typo cache
type TypoCache struct {
	Length        int    `yaml:"length"`
	CachingScheme string `yaml:"LFU"`
	WarmUp        bool   `yaml:"warmUp"`
}

// WaitList represents the wait list config
type WaitList struct {
	Length int `yaml:"length"`
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

	// check that token validity is set if not set to 15 min
	if s.Web.Reset.TokenValidity == 0*time.Second {
		log.Println("token validity is not set, using default value of 15 min")
		s.Web.Reset.TokenValidity = 15 * time.Minute
	}

	// check that rate-limiting is set if not set to 10
	if s.Web.Login.RateLimit == 0 {
		log.Println("rate limiting is not set, using default of 10")
		s.Web.Login.RateLimit = 10
	}

	// check that cookie validity is set if not set to 30 min
	if s.Web.Login.SessionValidity == 0*time.Second {
		log.Println("http session validity is not set, using default value of 30 min")
		s.Web.Login.SessionValidity = 30 * time.Minute
	}

	// validate typtop
	if s.Checker.TypTop != nil {
		s.validateTypTop()
	}

	return nil
}

func (s *Server) validateTypTop() {
	if s.Checker.TypTop.PublicKeyEncryption.KeyLength == 0 {
		log.Fatal("please set your pke key length")
	}

	if s.Checker.TypTop.PasswordBasedEncryption.KeyLength == 0 {
		log.Fatal("please set your pbe key length")
	}
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

	if s.Checker.TypTop != nil {
		numberOfCheckers++
	}

	if numberOfCheckers != 1 {
		return errors.New("more than one checker is defined in config - use only one of: always, blacklist, optimal")
	}

	return nil
}
