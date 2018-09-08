package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a graphql server projejct",
			Action: func(c *cli.Context) error {
				fmt.Println("not implemented")
				return nil
			},
		},
		{
			Name:    "gen",
			Aliases: []string{"g"},
			Usage:   "Generate models, resolvers and graphql server code",
			Action: func(c *cli.Context) error {

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
