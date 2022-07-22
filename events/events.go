package events

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/EfimoffN/authorBot/sqlapi"

	"github.com/google/uuid"
)

type IEvent interface {
	AddNewUser(userN string, userID int, chatID int64) error
	AddNewRefUserLink(userID int, link string) error
	RemoveRefUserLink(userID int, link string) error
	GetAllUserLinks(userID int) ([]*sqlapi.LinkRow, error)
	GetLinkByLink(link string) (*sqlapi.LinkRow, error)
}

type Event struct {
	SQLAPI *sqlapi.SQLAPI // исползовать интерфейс
	ctx    context.Context
}

var (
	ErrUserAlreadyAdded  = errors.New("such a user has already been added")
	ErrLinkAlreadyExists = errors.New("such a link already exists")
	ErrLinkNotDB         = errors.New("there is no such link in the database")
)

func NewBotEvents(sqlapi *sqlapi.SQLAPI, ctx context.Context) *Event {
	return &Event{
		SQLAPI: sqlapi,
		ctx:    ctx,
	}
}

func (api *Event) AddNewUser(userN string, userID int, chatID int64) error {
	user, err := api.getUser(userID)
	if err != nil {
		return e.Wrap("add new user failed with an error: ", err)
	}

	if user != nil {
		return e.Wrap("New user: ", ErrUserAlreadyAdded)
	}

	err = api.saveUser(userN, userID, chatID)
	if err != nil {
		return e.Wrap("sav new user failed with an error: ", err)
	}

	return nil
}

func (api *Event) AddNewRefUserLink(userID int, link string) error {
	linkRow, err := api.GetLinkByLink(link)
	if err != nil {
		return e.Wrap("get link by link failed with an error: ", err)
	}

	if linkRow == nil {
		linkRow.LinkID, err = api.saveLink(link)
		if err != nil {
			return e.Wrap("save link by link failed with an error: ", err)
		}
	}

	refRow, err := api.getRefByIDLinkUser(userID, link)
	if err != nil {
		return e.Wrap("get user link by id failed with an error: ", err)
	}

	if refRow != nil {
		return e.Wrap("New ref user link: ", ErrLinkAlreadyExists)
	}

	_, err = api.createRefUserLink(userID, link)
	if err != nil {
		return e.Wrap("create ref user link failed with an error: ", err)
	}

	return nil
}

func (api *Event) RemoveRefUserLink(userID int, link string) error {
	linkRow, err := api.GetLinkByLink(link)
	if err != nil {
		return e.Wrap("get link by link failed with an error: ", err)
	}

	if linkRow == nil {
		return e.Wrap("link row nil: ", ErrLinkNotDB)
	}

	err = api.deleteRefUserIDLinkID(userID, linkRow.LinkID)
	if err != nil {
		return e.Wrap("delete ref userID linkID: ", err)
	}

	return nil
}

func (api *Event) GetAllUserLinks(userID int) ([]*sqlapi.LinkRow, error) {
	refRow, err := api.SQLAPI.GetLinksUser(userID)
	if err != nil {
		return nil, e.Wrap("get all user links failed with an error: ", err)
	}

	return refRow, nil
}

func (api *Event) getRefByIDLinkUser(userID int, linkID string) (*sqlapi.RefRow, error) {
	refRow, err := api.SQLAPI.GetRefByIDLinkUser(userID, linkID)
	if err != nil {
		return nil, e.Wrap("get user link by id failed with an error: ", err)
	}

	return refRow, nil
}

func (api *Event) GetLinkByLink(link string) (*sqlapi.LinkRow, error) {
	linkRow, err := api.SQLAPI.GetLinkByLink(link)
	if err != nil {
		return nil, e.Wrap("get link by link failed with an error: ", err)
	}

	return linkRow, nil
}

func (api *Event) getUser(userID int) (*sqlapi.UserRow, error) {
	userRow, err := api.SQLAPI.GetUserByID(userID)
	if err != nil {
		return nil, e.Wrap("get user by ID failed with an error: ", err)
	}

	return userRow, nil
}

func (api *Event) getRefLinksUser(userID int) ([]*sqlapi.RefRow, error) {
	refRow, err := api.SQLAPI.GetRefLinksUser(userID)
	if err != nil {
		return nil, e.Wrap("get all ref user links failed with an error: ", err)
	}

	return refRow, nil
}

func (api *Event) createRefUserLink(userID int, linkID string) (string, error) {
	uuid := getUUID()
	err := api.SQLAPI.AddRefLinkUser(api.ctx, uuid, linkID, userID)
	if err != nil {
		return "", e.Wrap("create new ref user link failed with an error: ", err)
	}

	return uuid, nil
}

func (api *Event) saveUser(userN string, userID int, chatID int64) error {
	err := api.SQLAPI.AddUser(api.ctx, userN, userID, chatID)
	if err != nil {
		return e.Wrap("save new user failed with an error: ", err)
	}

	return nil
}

func (api *Event) saveLink(link string) (string, error) {
	uuid := getUUID()
	err := api.SQLAPI.AddLink(api.ctx, link, uuid)
	if err != nil {
		return "", e.Wrap("save new link failed with an error: ", err)
	}

	return uuid, nil
}

func (api *Event) deleteRefUserIDLinkID(userID int, linkID string) error {
	err := api.SQLAPI.RemoveRefByUserIDLinkID(userID, linkID)
	if err != nil {
		return e.Wrap("remove ref by userid linkid failed with an error: ", err)
	}

	return nil
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
