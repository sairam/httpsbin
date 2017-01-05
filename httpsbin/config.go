package httpsbin

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
)

// AppConfig configuration to start the server
type AppConfig struct {
	DataDir     string `toml:"data_dir"`
	Persistence string `toml:"persistence"`

	ServerProto string `toml:"server_proto"`
	ServerHost  string `toml:"server_host"`
	LocalServer string `toml:"server_local"`

	CleanupStrategy   string
	CleanupAfter      int
	MaxFilesToDisplay int

	CacheView bool `toml:"cache"`

	once sync.Once
}

// Config used outside
var Config = &AppConfig{}

// InitConfig should be called to initialise the config
func InitConfig(file string) {
	Config.once.Do(func() { Config.init(file) })
}

func (config *AppConfig) init(file string) {
	Config.CleanupStrategy = "idle"
	Config.CleanupAfter = 6 * 60 * 60 // in seconds, 6 hours
	Config.MaxFilesToDisplay = 20

	if file != "" && config.loadFromFile(file) {
		log.Println("Loaded from Config")
		return
	}

	config.DataDir = os.Getenv("DATA_DIR")        // "data/"
	config.Persistence = os.Getenv("PERSISTENCE") // "filesystem" / optional "inmemory". DataDir is ignored if inmemory

	Config.LocalServer = os.Getenv("LOCALSERVER")  // "localhost:3000"
	Config.ServerProto = os.Getenv("SERVER_PROTO") // "http"
	Config.ServerHost = os.Getenv("SERVER_HOST")   // "localhost:3000"

	Config.CacheView = os.Getenv("CACHE_VIEW") == "true" // false means not to cache, true means to cache

	log.Println("Loaded from env variables")
}

func (config *AppConfig) loadFromFile(file string) bool {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return false
	}
	toml.Unmarshal(data, &config)
	return true
}
