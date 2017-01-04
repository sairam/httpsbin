package main

import (
	"errors"

	"github.com/spf13/afero"
)

const (
	persistenceFs  = "filesystem"
	persistenceMem = "inmemory"
)

var fsutil *afero.Afero

// InitPersist creates directory if it does not exist
func InitPersist() {
	var appFs afero.Fs
	if Config.Persistence == persistenceFs {
		appFs = afero.NewBasePathFs(afero.NewOsFs(), Config.DataDir)
	} else {
		appFs = afero.NewMemMapFs()
	}
	fsutil = &afero.Afero{Fs: appFs}
}

func readFile(filename string) ([]byte, error) {
	return fsutil.ReadFile(filename)
}

// CleanStaleFiles cleans stale files / directories periodically
func CleanStaleFiles() {
	// Config.CleanupStrategy == 'idle'
	// Config.CleanupAfter , delete directories after X time
}

// CleanUpMaxItemsInDir is called via a go()
func CleanUpMaxItemsInDir(dir string) {

}

func createNewDir(dir string) (string, error) {
	ok, _ := ifDirExists(dir)
	if ok {
		return "", errors.New("Directory already exists")
	}
	fsutil.Mkdir(dir, 0700)
	return dir, nil
}

func ifDirExists(dir string) (bool, error) {
	_, err := fsutil.Stat(dir)
	if err != nil {
		return false, err
	}
	return fsutil.IsDir(dir)
}
