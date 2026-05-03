package main

import (
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/app"
)

func main() {
	start := app.NewApp()

	start.Run()
}
