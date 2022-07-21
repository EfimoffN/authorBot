package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Storage interface {
}

// Commands ...
type Commands struct {
	StartBot     string
	StopBot      string
	StartStopBot string
	CreateCode   string
	EditCode     string
	RegisterUser string
	StartPigeon  string
	AddNameBot   string
	EditNameBot  string
}

type BotSvc struct {
	storage  Storage
	commands Commands
}

func NewBotSvc(s Storage, commands Commands) *BotSvc {
	return &BotSvc{
		storage:  s,
		commands: commands,
	}
}

func (b *BotSvc) ProcessingComands(message *tgbotapi.Message, bot *tgbotapi.BotAPI) error {

}
