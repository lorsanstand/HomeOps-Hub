package main

import "github.com/lorsanstand/HomeOps-Hub/internal/agent/app"

func main() {
	start := app.NewApp()
	start.Run()
}
