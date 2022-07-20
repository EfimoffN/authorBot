package sqlapi

// UserRow ...
type UserRow struct {
	UserID   string `db:"userid"`
	NameUser string `db:"nameuser"`
	ChatID   string `db:"chatid"`
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
	UserID string `db:"userid"`
}
