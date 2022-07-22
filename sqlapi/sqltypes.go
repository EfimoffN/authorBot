package sqlapi

// UserRow ...
type UserRow struct {
	UserID   int    `db:"userid"`
	NameUser string `db:"nameuser"`
	ChatID   int64  `db:"chatid"`
}

// LinkRow ...
type LinkRow struct {
	LinkID string `db:"linkid"`
	Link   string `db:"link"`
}

// RefRow ...
type RefRow struct {
	RefID  string `db:"refid"`
	LinkID string `db:"linkid"`
	UserID int    `db:"userid"`
}
