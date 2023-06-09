package main

import (
	"github.com/ystyle/jvms/utils/cmd"
	"log"
	"os"
)

func main() {
	app := cmd.App()
	if err := app.Run(os.Args); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
