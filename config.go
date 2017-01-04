package main

import "sync"

// AppConfig configuration to start the server
type AppConfig struct {
	DataDir           string
	Persistence       string
	LocalServer       string
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
	config.DataDir = "data/"
	config.Persistence = "filesystem" // optional "inmemory". DataDir is ignored if inmemory
	Config.LocalServer = "localhost:3000"
	Config.CleanupStrategy = "idle"
	Config.CleanupAfter = 6 * 60 * 60 // in seconds, 6 hours
	Config.MaxFilesToDisplay = 20
}
