package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beinan/gql-server/codegen"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gqltools"
	config := codegen.GenConfig{
		SchemaPath: "./schema",
		GenPath:    "./gen",
	}
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
			Subcommands: []cli.Command{
				{
					Name:    "model",
					Aliases: []string{"m"},
					Usage:   "generate models",
					Action: func(c *cli.Context) error {
						codegen.GenerateModel(config, os.Stdout)
						return nil
					},
				},
				{
					Name:    "resolver",
					Aliases: []string{"r"},
					Usage:   "generate resolvers",
					Action: func(c *cli.Context) error {
						codegen.GenerateResolver(config, os.Stdout)
						return nil
					},
				},
				{
					Name:    "gqlresolver",
					Aliases: []string{"gr"},
					Usage:   "generate gql resolvers",
					Action: func(c *cli.Context) error {
						codegen.GenerateGqlResolver(config, os.Stdout)
						return nil
					},
				},
				{
					Name:  "mock",
					Usage: "generate mocks",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
			}},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
