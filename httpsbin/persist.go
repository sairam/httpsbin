package httpsbin

import (
	"errors"
	"log"
	"sort"
	"strconv"

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

// FileItem ...
type FileItem struct {
	Name      string
	Timestamp int64
}

// ParseTime parses time from filename
func (fi *FileItem) ParseTime() (err error) {
	fi.Timestamp, err = strconv.ParseInt(fi.Name, 10, 64)
	return
}

// FileList ..
type FileList []*FileItem

func (a FileList) Len() int           { return len(a) }
func (a FileList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FileList) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

// CleanUpMaxItemsInDir is called via a go()
func CleanUpMaxItemsInDir(dir string) {
	fis, err := fsutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	if len(fis) <= Config.MaxFilesToDisplay {
		// not gonna try to cleanup
		return
	}

	filelist := make([]*FileItem, 0, len(fis))
	for _, fi := range fis {
		fileitem := &FileItem{fi.Name(), 0}
		if err := fileitem.ParseTime(); err == nil {
			filelist = append(filelist, fileitem)
		} else {
			log.Println(err)
		}
	}
	sort.Sort(sort.Reverse(FileList(filelist)))

	if len(filelist) > Config.MaxFilesToDisplay {
		for _, ir := range filelist[Config.MaxFilesToDisplay:] {
			fsutil.Remove(MergeOSPath(dir, ir.Name))
		}
	}

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
