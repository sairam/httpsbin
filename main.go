package main

// import _ "github.com/sairam/httpsbin"
func main() {
	InitConfig()
	InitPersist()

	go CleanStaleFiles()
	InitRouter()
}
