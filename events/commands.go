package events

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/EfimoffN/authorBot/sqlapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

type Commands struct {
	BotAPI *tgbotapi.BotAPI
	SQLAPI *sqlapi.SQLAPI // испольщовать интерфейс
	ctx    context.Context
}

var (
	IsNotLink          = errors.New("text is not link")
	ErrUnknownMetaType = errors.New("unknown meta type")
	UserAlreadyAdded   = errors.New("such a user has already been added")
)

func NewBotCommands(botAPI *tgbotapi.BotAPI, sqlapi *sqlapi.SQLAPI, ctx context.Context) *Commands {
	return &Commands{
		BotAPI: botAPI,
		SQLAPI: sqlapi,
		ctx:    ctx,
	}
}

func (api *Commands) AddNewUser(userN, userID, chatID string) error {
	user, err := api.getUser(userID)
	if err != nil {
		return e.Wrap("add new user failed with an error: ", err)
	}

	if user != nil {
		return e.Wrap("New user: ", UserAlreadyAdded)
	}

	err = api.saveUser(userN, userID, chatID)
	if err != nil {
		return e.Wrap("sav new user failed with an error: ", err)
	}

	return nil
}

func (api *Commands) AddNewRefUserLink(userID, link string) error {
	linkRow, err := api.getLinkByLink(link)
	if err != nil {
		return e.Wrap("get link by link failed with an error: ", err)
	}

	if linkRow == nil {
		linkRow.LinkID, err = api.saveLink(link)
		if err != nil {
			return e.Wrap("save link by link failed with an error: ", err)
		}
	}

	return nil
}

func (api *Commands) getLinksUser(userID string) ([]*sqlapi.LinkRow, error) {
	refRow, err := api.SQLAPI.GetLinksUser(userID)
	if err != nil {
		return nil, e.Wrap("get all user links failed with an error: ", err)
	}

	return refRow, nil
}

func (api *Commands) getLinkByLink(link string) (*sqlapi.LinkRow, error) {
	linkRow, err := api.SQLAPI.GetLinkByLink(link)
	if err != nil {
		return nil, e.Wrap("get link by link failed with an error: ", err)
	}

	return linkRow, nil
}

func (api *Commands) getUser(userID string) (*sqlapi.UserRow, error) {
	userRow, err := api.SQLAPI.GetUserByID(userID)
	if err != nil {
		return nil, e.Wrap("get user by ID failed with an error: ", err)
	}

	return userRow, nil
}

func (api *Commands) getRefLinksUser(userID string) ([]*sqlapi.RefRow, error) {
	refRow, err := api.SQLAPI.GetRefLinksUser(userID)
	if err != nil {
		return nil, e.Wrap("get all ref user links failed with an error: ", err)
	}

	return refRow, nil
}

func (api *Commands) createRefUserLink(userID, linkID string) (string, error) {
	uuid := getUUID()
	err := api.SQLAPI.AddRefLinkUser(api.ctx, uuid, linkID, userID)
	if err != nil {
		return "", e.Wrap("create new ref user link failed with an error: ", err)
	}

	return uuid, nil
}

func (api *Commands) saveUser(userN, userID, chatID string) error {
	err := api.SQLAPI.AddUser(api.ctx, userN, userID, chatID)
	if err != nil {
		return e.Wrap("save new user failed with an error: ", err)
	}

	return nil
}

func (api *Commands) saveLink(link string) (string, error) {
	uuid := getUUID()
	err := api.SQLAPI.AddLink(api.ctx, link, uuid)
	if err != nil {
		return "", e.Wrap("save new link failed with an error: ", err)
	}

	return uuid, nil
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func getUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}
