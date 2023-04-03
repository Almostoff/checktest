package main

import (
	_ "github.com/lib/pq"
	app "wbL0/server_sub/app"
	config "wbL0/server_sub/config"
)

func main() {
	config := new(config.Config)
	config.InitFile()
	app := app.InitApp(*config)
	app.Run()
}
