package commands

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/EfimoffN/authorBot/events"
	"github.com/EfimoffN/authorBot/lib/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	RmvCmd     = "/rmv"
	HelpCmd    = "/help"
	StartCmd   = "/start"
	AllLinkCmd = "/all"
)

type ICommands interface {
	NewBotCommands(botAPI *tgbotapi.BotAPI, event events.IEvent, ctx context.Context) *Commands
}

type Commands struct {
	BotAPI *tgbotapi.BotAPI
	Event  events.IEvent
	CTX    context.Context
}

func NewBotCommands(botAPI *tgbotapi.BotAPI, event events.IEvent, ctx context.Context) *Commands {
	return &Commands{
		BotAPI: botAPI,
		Event:  event,
		CTX:    ctx,
	}
}

func (c *Commands) DoCommand(message *tgbotapi.Message) error {
	text := strings.TrimSpace(message.Text)

	if isAddCmd(text) {
		return c.saveLink(message)
	}

	if isRemoveLink(text) {
		return c.removeLink(message)
	}

	switch text {
	case StartCmd:
		return c.sendStart(message)
	case HelpCmd:
		return c.sendHelp(message)
	case AllLinkCmd:
		return c.getAllLinks(message)
	default:
		return c.sendUnknownCommand(message)
	}

}

func (c *Commands) sendHelp(message *tgbotapi.Message) error {
	m := tgbotapi.NewMessage(message.Chat.ID, msgHelp)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err := c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func (c *Commands) sendStart(message *tgbotapi.Message) error {
	err := c.addNewUser(message)
	if err != nil {
		return e.Wrap("save new user failed with an error: ", err)
	}

	m := tgbotapi.NewMessage(message.Chat.ID, msgHello)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err = c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func (c *Commands) sendUnknownCommand(message *tgbotapi.Message) error {
	m := tgbotapi.NewMessage(message.Chat.ID, msgUnknownCommand)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err := c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func (c *Commands) addNewUser(message *tgbotapi.Message) error {
	err := c.Event.AddNewUser(message.From.UserName, message.From.ID, message.Chat.ID)

	if err != nil && !errors.Is(err, events.ErrUserAlreadyAdded) {
		return e.Wrap("save new user failed with an error: ", err)
	}

	return nil
}

func (c *Commands) saveLink(message *tgbotapi.Message) error {
	err := c.Event.AddNewRefUserLink(message.From.ID, message.Text)
	if err != nil && !errors.Is(err, events.ErrLinkAlreadyExists) {
		return e.Wrap("save link failed with an error: ", err)
	}
	msg := "Ссылка сохранена для отслеживания."

	if errors.Is(err, events.ErrLinkAlreadyExists) {
		msg = "Такая ссылка уже добавлялась."
	}

	m := tgbotapi.NewMessage(message.Chat.ID, msg)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err = c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func (c *Commands) removeLink(message *tgbotapi.Message) error {
	text := strings.TrimSpace(message.Text)[2:]

	link := strings.TrimSpace(text)

	err := c.Event.RemoveRefUserLink(message.From.ID, link)
	if err != nil {
		return e.Wrap("remove link failed with an error: ", err)
	}

	msg := "Ссылка больше не отслеживается."

	m := tgbotapi.NewMessage(message.Chat.ID, msg)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err = c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func (c *Commands) getAllLinks(message *tgbotapi.Message) error {
	linkRows, err := c.Event.GetAllUserLinks(message.From.ID)
	if err != nil {
		return e.Wrap("get links failed with an error: ", err)
	}

	msg := "У вас нет отслеживаемых ссылок"

	if len(linkRows) > 0 {
		msg = "Ваши отслеживаемые ссылки:\n"

		for _, l := range linkRows {
			msg = msg + l.Link + "\n"
		}
	}

	m := tgbotapi.NewMessage(message.Chat.ID, msg)
	replyKeyboardHide := tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	m.ReplyMarkup = replyKeyboardHide
	_, err = c.BotAPI.Send(m)
	if err != nil {
		return e.Wrap("Sending the message failed with an error: ", err)
	}

	return nil
}

func isAddCmd(text string) bool {
	return isURL(text)
}

// TODO проверять ссылку, работает или нет
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func isRemoveLink(text string) bool {
	prefix := strings.TrimSpace(text)[0:2]
	remove := false

	if prefix == "rm" {
		remove = true
	}

	if remove {
		link := strings.TrimSpace(text[2:])

		remove = isURL(link)
	}

	return remove
}
