package main

import "github.com/lorsanstand/HomeOps-Hub/internal/agent/app"

func main() {
	start, err := app.NewApp()
	if err != nil {
		return
	}
	start.Run()
}
