package service

import (
	"log"

	"github.com/EfimoffN/authorBot/commands"
	"github.com/EfimoffN/authorBot/lib/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(cmd *commands.Commands, bot *tgbotapi.BotAPI, timeout int) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return e.Wrap("Get updates chan", err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if err := cmd.DoCommand(update.Message); err != nil {
			log.Println("Processing commands: ", err.Error())
		}
	}

	return nil
}
