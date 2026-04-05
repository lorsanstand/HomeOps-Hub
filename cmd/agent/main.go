package main

import "github.com/lorsanstand/HomeOps-Hub/internal/hub/app"

func main() {
	start := app.NewApp()
	start.Run()
}
