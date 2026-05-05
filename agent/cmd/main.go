package main

import (
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/app"
)

func main() {
	start, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	start.Run()
}
