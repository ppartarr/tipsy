package session

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	Name      = "session"
	keyUserID = "userID"
)

var (
	// Store contains all sessions
	Store      *sessions.FilesystemStore
	db         bleve.Index
	StorageDir string
	secret     string

	Log = logrus.New()
)

func init() {
	secret = randomID()
	Log.Formatter = &prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
	}
}

// InitStore initializes the session store
// secret should be generated using securecookie.GenerateRandomKey()
// maxAge is supplied in seconds
// pass maxAge=0 for no timeout
func InitStore(secret string, dbHandle bleve.Index, maxAge int) {
	sessionMaxAge = int64(maxAge)
	err := os.MkdirAll(filepath.Join(StorageDir, "sessions"), 0o755)
	if err != nil && err != os.ErrExist {
		Log.WithError(err).Warn("failed to create sessions directory")
	}

	// seed rand pkg with startup timestamp
	rand.Seed(time.Now().UnixNano())

	Store = sessions.NewFilesystemStore(filepath.Join(StorageDir, "sessions"), []byte(secret))
	db = dbHandle

	Store.MaxAge(maxAge)
}
