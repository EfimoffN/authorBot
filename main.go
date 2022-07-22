package main

import (
	"context"
	"log"
	"os"

	"github.com/EfimoffN/authorBot/commands"
	"github.com/EfimoffN/authorBot/config"
	"github.com/EfimoffN/authorBot/events"
	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/EfimoffN/authorBot/sqlapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
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
			ctx := context.Background()

			cfg, err := config.CreateConfig(cCtx.String("config"))
			if err != nil {
				return e.Wrap("Create config", err)
			}

			db, err := connectDB(cfg.ConnectPostgres)
			if err != nil {
				return e.Wrap("connect DB", err)
			}
			defer db.Close()

			sqlAPI := sqlapi.NewSQLAPI(db)

			event := events.NewBotEvents(sqlAPI, ctx)

			bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
			if err != nil {
				return e.Wrap("New bot API", err)
			}

			bot.Debug = true

			cmnd := commands.NewBotCommands(bot, event, ctx)

			//TODO вынести в отдельный пакет
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60

			updates, err := bot.GetUpdatesChan(u)
			if err != nil {
				return e.Wrap("Get updates chan", err)
			}
			for update := range updates {
				if update.Message == nil { // ignore any non-Message Updates
					continue
				}

				if err := cmnd.DoCommand(update.Message); err != nil {
					log.Println("Processing commands: ", err.Error())
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// connectDB ...
func connectDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		log.Println("sqlx.Open failed with an error: ", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("DB.Ping failed with an error: ", err.Error())
		return nil, err
	}

	return db, err
}
