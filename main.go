package main

import (
	"github.com/azinudinachzab/scr-syky-tech-test/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app.New().Run()
}
