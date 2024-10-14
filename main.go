package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "ossie",
		Usage: "A powerful Tool to manage Openstack Contexts",
		Action: func(*cli.Context) error {
			fmt.Println("ossie says hi!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
