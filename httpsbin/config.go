package httpsbin

import (
	"os"
	"sync"
)

// AppConfig configuration to start the server
type AppConfig struct {
	DataDir     string
	Persistence string

	ServerProto string
	ServerHost  string
	LocalServer string

	CleanupStrategy   string
	CleanupAfter      int
	MaxFilesToDisplay int

	once sync.Once
}

// Config used outside
var Config = &AppConfig{}

// InitConfig should be called to initialise the config
func InitConfig() {
	Config.once.Do(func() { Config.init() })
}

func (config *AppConfig) init() {
	config.DataDir = os.Getenv("DATA_DIR")        // "data/"
	config.Persistence = os.Getenv("PERSISTENCE") // "filesystem" / optional "inmemory". DataDir is ignored if inmemory

	Config.LocalServer = os.Getenv("LOCALSERVER")  // "localhost:3000"
	Config.ServerProto = os.Getenv("SERVER_PROTO") // "http"
	Config.ServerHost = os.Getenv("SERVER_HOST")   // "localhost:3000"

	Config.CleanupStrategy = "idle"
	Config.CleanupAfter = 6 * 60 * 60 // in seconds, 6 hours
	Config.MaxFilesToDisplay = 20
}
