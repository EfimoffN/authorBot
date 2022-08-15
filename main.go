package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EfimoffN/authorBot/commands"
	"github.com/EfimoffN/authorBot/config"
	"github.com/EfimoffN/authorBot/events"
	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/EfimoffN/authorBot/service"
	"github.com/EfimoffN/authorBot/sqlapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli/v2"
)

func main() {

	pgPass := os.Getenv("pgPass")
	fmt.Println(pgPass)

	return

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
			ctx := context.Background()

			cfg, err := config.CreateConfig(cCtx.String("config"))
			if err != nil {
				return e.Wrap("Create config: ", err)
			}

			db, err := sqlapi.ConnectDB(cfg.ConnectPostgres)
			if err != nil {
				return e.Wrap("connect DB: ", err)
			}
			defer db.Close()

			sqlAPI := sqlapi.NewSQLAPI(db)

			event := events.NewBotEvents(sqlAPI, ctx)

			bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
			if err != nil {
				return e.Wrap("New bot API: ", err)
			}

			bot.Debug = true

			cmd := commands.NewBotCommands(bot, event, ctx)

			err = service.Start(cmd, bot, cfg.Timeout)
			if err != nil {
				return e.Wrap("service start: ", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
