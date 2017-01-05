package main

import (
	"os"

	bin "./httpsbin"
)

func main() {
	var configfile string
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	}

	bin.InitConfig(configfile)
	bin.InitPersist()

	go bin.CleanStaleFiles()
	bin.InitView()
	bin.InitRouter()
}
