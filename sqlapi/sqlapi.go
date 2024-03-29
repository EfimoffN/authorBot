package sqlapi

import (
	"context"

	"github.com/EfimoffN/authorBot/lib/e"
	"github.com/jmoiron/sqlx"
)

type ISQLAPI interface {
	GetLinksUser(userID int) ([]*LinkRow, error)
	GetRefByIDLinkUser(userID int, linkID string) (*RefRow, error)
	GetUserByID(userID int) (*UserRow, error)
	GetRefLinksUser(userID int) ([]*RefRow, error)
	AddRefLinkUser(ctx context.Context, refID string, linkID string, userID int) error
	AddUser(ctx context.Context, userN string, userID int, chatID int64) error
	RemoveRefByUserIDLinkID(userID int, linkID string) error
	GetLinkByLink(lnk string) (*LinkRow, error)
	AddLink(ctx context.Context, link, linkID string) error
}

type SQLAPI struct {
	db *sqlx.DB
}

func NewSQLAPI(db *sqlx.DB) *SQLAPI {
	return &SQLAPI{
		db: db,
	}
}

func ConnectDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, e.Wrap("sqlx.Open failed with an error: ", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("DB.Ping failed with an error: ", err)
	}

	return db, err
}

func (api *SQLAPI) GetUserByID(userID int) (*UserRow, error) {
	userRow := []UserRow{}

	err := api.db.Select(&userRow, "SELECT * FROM prj_user WHERE userid = $1;", userID)
	if err != nil {
		return nil, e.Wrap("GetUserByID api.db.Select failed with an error: ", err)
	}

	if len(userRow) == 1 {
		return &userRow[0], nil
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

func (api *SQLAPI) GetRefLinksUser(userID int) ([]*RefRow, error) {
	refRow := []*RefRow{}

	err := api.db.Select(&refRow, "SELECT * FROM ref_link_user WHERE userid = $1;", userID)
	if err != nil {
		return nil, e.Wrap("GetLinkByLink api.db.Select failed with an error: ", err)
	}

	return refRow, err
}

func (api *SQLAPI) GetLinksUser(userID int) ([]*LinkRow, error) {
	linkRow := []*LinkRow{}

	err := api.db.Select(&linkRow, "SELECT prj_link.linkid, prj_link.link FROM ref_link_user JOIN prj_link ON prj_link.linkid = ref_link_user.linkid WHERE ref_link_user.userid = $1;", userID)
	if err != nil {
		return nil, e.Wrap("GetLinksUser api.db.Select failed with an error: ", err)
	}

	return linkRow, err
}

func (api *SQLAPI) GetRefByIDLinkUser(userID int, linkID string) (*RefRow, error) {
	refRow := []RefRow{}

	err := api.db.Select(&refRow, "SELECT * FROM ref_link_user WHERE userid = $1 AND linkid = $2;", userID, linkID)
	if err != nil {
		return nil, e.Wrap("GetRefByIDLinkUser api.db.Select failed with an error: ", err)
	}

	if len(refRow) == 1 {
		return &refRow[0], nil
	}

	return nil, err
}

func (api *SQLAPI) AddUser(ctx context.Context, userN string, userID int, chatID int64) error {
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

func (api *SQLAPI) AddRefLinkUser(ctx context.Context, refID string, linkID string, userID int) error {
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

func (api *SQLAPI) RemoveRefByUserIDLinkID(userID int, linkID string) error {
	_, err := api.db.Exec("DELETE FROM ref_link_user WHERE userID = $1 AND linkid = $2;", userID, linkID)
	if err != nil {
		return e.Wrap("DELETE row ref link by userid adn linkid failed with an error: ", err)
	}

	return nil
}

func (api *SQLAPI) RemoveUser(userID int) error {
	_, err := api.db.Exec("DELETE FROM prj_user WHERE userid = $1;", userID)
	if err != nil {
		return e.Wrap("DELETE row user failed with an error: ", err)
	}

	return nil
}
