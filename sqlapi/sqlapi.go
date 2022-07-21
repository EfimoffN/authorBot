package sqlapi

import (
	"context"
	"strings"

	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLAPI struct {
	db *sqlx.DB
}

func NewSQLAPI(db *sqlx.DB) *SQLAPI {
	return &SQLAPI{
		db: db,
	}
}

func (api *SQLAPI) GetUserByID(userID string) (*UserRow, error) {
	userRow := []UserRow{}

	err := api.db.Select(&userRow, "SELECT * FROM prj_user WHERE userid = $1;", userID)
	if err != nil {
		return nil, e.Wrap("GetUserByID api.db.Select failed with an error: ", err)
	}

	if len(userRow) == 1 {
		return &userRow[0], err
	}

	return nil, err
}

func (api *SQLAPI) GetLinkByLink(lnk string) (*LinkRow, error) {
	linkRow := []LinkRow{}

	err := api.db.Select(&linkRow, "SELECT * FROM prj_link WHERE link = $1;", lnk)
	if err != nil {
		return nil, e.Wrap("GetLinkByLink api.db.Select failed with an error: ", err)
	}

	if len(linkRow) == 1 {
		return &linkRow[0], err
	}

	return nil, err
}

func (api *SQLAPI) GetLinksUser(userID string) ([]*RefRow, error) {
	refRow := []*RefRow{}

	err := api.db.Select(&refRow, "SELECT * FROM ref_link_user WHERE userid = $1;", userID)
	if err != nil {
		return nil, e.Wrap("GetLinkByLink api.db.Select failed with an error: ", err)
	}

	return refRow, err
}

func (api *SQLAPI) AddUser(ctx context.Context, userN, userID, chatID string) error {
	const query = `INSERT INTO prj_user(userid, nameuser, chatid) VALUES (:userid, :nameuser, :chatid) ON CONFLICT DO NOTHING;`

	user := UserRow{
		UserID:   userID,
		NameUser: userN,
		ChatID:   chatID,
	}

	if _, err := api.db.NamedExecContext(ctx, query, user); err != nil {
		return e.Wrap("INSERT user failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) AddLink(ctx context.Context, link, linkID string) error {
	const query = `INSERT INTO prj_link(linkid, link) VALUES (:linkid, :link) ON CONFLICT DO NOTHING;`

	linkR := LinkRow{
		LinkID: linkID,
		Link:   link,
	}

	if _, err := api.db.NamedExecContext(ctx, query, linkR); err != nil {
		return e.Wrap("INSERT link failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) AddRefLinkUser(ctx context.Context, refID, linkID, userID string) error {
	const query = `INSERT INTO ref_link_user(refid, linkid, userid) VALUES (:refid, :linkid, :userid) ON CONFLICT DO NOTHING;`

	refR := RefRow{
		RefID:  refID,
		LinkID: linkID,
		UserID: userID,
	}

	if _, err := api.db.NamedExecContext(ctx, query, refR); err != nil {
		return e.Wrap("INSERT ref_link_user failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) RemoveRefLinkUser(refID string) error {
	_, err := api.db.Exec("DELETE FROM ref_link_user WHERE refid = $1;", refID)
	if err != nil {
		return e.Wrap("DELETE row ref link failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) RemoveUser(userID string) error {
	_, err := api.db.Exec("DELETE FROM prj_user WHERE userid = $1;", userID)
	if err != nil {
		return e.Wrap("DELETE row user failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) RemoveLinksUser(userID string) error {
	_, err := api.db.Exec("DELETE FROM ref_link_user WHERE userid = $1;", userID)
	if err != nil {
		return e.Wrap("DELETE all user links failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) GetUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}
