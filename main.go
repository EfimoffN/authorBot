package main

import (
	"log"
	"os"

	"github.com/EfimoffN/authorBot/config"
	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name: "authorBot",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration",
			},
		},
		Action: func(cCtx *cli.Context) error {

			cfg, err := config.CreateConfig(cCtx.String("config"))
			if err != nil {
				return e.Wrap("create config", err)
			}

			// connect to db

			// eventProc

			// consumer

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
