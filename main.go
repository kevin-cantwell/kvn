package main

import (
	"os"

	"github.com/kevin-cantwell/kvn/cmd/admin"
	"github.com/kevin-cantwell/kvn/cmd/web"
)

func main() {
	switch app := os.Args[1]; app {
	case "web":
		web.Run()
	case "admin":
		admin.Run()
	default:
		panic("unknown app: " + app)
	}
}
