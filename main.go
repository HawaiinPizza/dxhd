package main

import (
	"github.com/dakyskye/dxhd/app"
	"github.com/dakyskye/dxhd/logger"
)

func main() {
	a, err := app.Init()
	if err != nil {
		logger.L().WithError(err).Fatalln("can not initialise the app")
	}

	err = a.Parse()
	if err != nil {
		logger.L().WithError(err).Fatalln("something went wrong while parsing the app arguments")
	}

	err = a.Start()
	if err != nil {
		logger.L().WithError(err).Fatal("something went wrong while running the app")
	}
}
