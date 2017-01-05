package main

import bin "./httpsbin"

func main() {
	bin.InitConfig()
	bin.InitPersist()

	go bin.CleanStaleFiles()
	bin.InitView()
	bin.InitRouter()
}
