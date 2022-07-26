package sqlapi

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_GetUserByID(t *testing.T) {
	columns := []string{"userid", "nameuser", "chatid"}

	const expectedQuery = `SELECT (.+) FROM prj_user WHERE userid = (.+);`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "get user",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,1"))
			},
			wantErr: false,
		},
		{
			name: "get user nil",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,0,0"))
			},
			wantErr: false,
		},
		{
			name: "get user error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,1"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)

			_, err = api.GetUserByID(123)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("SetNewUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetLinkByLink(t *testing.T) {
	columns := []string{"linkid", "link"}

	const expectedQuery = "SELECT (.+) FROM prj_link WHERE link = (.+);"

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "get link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs("http//test.test").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1"))
			},
			wantErr: false,
		},
		{
			name: "get link nil",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs("http//test.test").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,0"))
			},
			wantErr: false,
		},
		{
			name: "get link error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs("http//test.test").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)

			_, err = api.GetLinkByLink("http//test.test")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkByLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetLinkByLink() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetRefLinksUser(t *testing.T) {
	columns := []string{"refid", "userid", "linkid"}

	const expectedQuery = "SELECT (.+) FROM ref_link_user WHERE userid = (.+);"

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "get link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,1"))
			},
			wantErr: false,
		},
		{
			name: "get link nil",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,0,0"))
			},
			wantErr: false,
		},
		{
			name: "get link error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)

			_, err = api.GetRefLinksUser(123)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefLinksUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetRefLinksUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetRefByIDLinkUser(t *testing.T) {
	columns := []string{"refid", "userid", "linkid"}

	const expectedQuery = "SELECT (.+) FROM ref_link_user WHERE userid = (.+) AND linkid = (.+);"

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "get link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123, "link").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1,1"))
			},
			wantErr: false,
		},
		{
			name: "get link nil",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123, "link").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,0,0"))
			},
			wantErr: false,
		},
		{
			name: "get link error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123, "link").
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)

			_, err = api.GetRefByIDLinkUser(123, "link")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefByIDLinkUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetRefByIDLinkUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetLinksUser(t *testing.T) {
	columns := []string{"refid", "userid", "linkid"}

	const expectedQuery = "SELECT prj_link.linkid, prj_link.link FROM ref_link_user JOIN prj_link ON prj_link.linkid = ref_link_user.linkid WHERE ref_link_user.userid = (.+);"

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "get link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1"))
			},
			wantErr: false,
		},
		{
			name: "get link nil",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,0,0"))
			},
			wantErr: false,
		},
		{
			name: "get link error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows(columns).FromCSVString("0,1"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)

			_, err = api.GetLinksUser(123)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinksUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetLinksUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_AddUser(t *testing.T) {
	ctx := context.Background()
	user := UserRow{
		UserID:   123,
		NameUser: "userN",
		ChatID:   321,
	}

	const expectedQuery = `INSERT INTO prj_user\(userid, nameuser, chatid\)	VALUES (.+)	ON CONFLICT DO NOTHING;`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success add new user",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(user.UserID, user.NameUser, user.ChatID).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: false,
		},
		{
			name: "error on add new user",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(user.UserID, user.NameUser, user.ChatID).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)
			if err := api.AddUser(ctx, user.NameUser, user.UserID, user.ChatID); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("AddUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_AddLink(t *testing.T) {
	ctx := context.Background()
	linkR := LinkRow{
		LinkID: "linkID",
		Link:   "link",
	}
	const expectedQuery = `INSERT INTO prj_link\(linkid, link\) VALUES (.+) ON CONFLICT DO NOTHING;`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success add new link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(linkR.LinkID, linkR.Link).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: false,
		},
		{
			name: "error on add new link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(linkR.LinkID, linkR.Link).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)
			if err := api.AddLink(ctx, linkR.Link, linkR.LinkID); (err != nil) != tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("AddLink() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_AddRefLinkUser(t *testing.T) {
	ctx := context.Background()
	refR := RefRow{
		RefID:  "refID",
		LinkID: "linkID",
		UserID: 123,
	}

	const expectedQuery = `INSERT INTO ref_link_user\(refid, linkid, userid\) VALUES (.+) ON CONFLICT DO NOTHING;`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success add new link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(refR.RefID, refR.LinkID, refR.UserID).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: false,
		},
		{
			name: "error on add new link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(refR.RefID, refR.LinkID, refR.UserID).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)
			if err := api.AddRefLinkUser(ctx, refR.RefID, refR.LinkID, refR.UserID); (err != nil) != tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("AddLink() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_RemoveUser(t *testing.T) {
	const expectedQuery = `DELETE FROM prj_user WHERE userid = (.+);`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "delete link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(123).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: false,
		},
		{
			name: "delete link err",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(123).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)
			if err := api.RemoveUser(123); (err != nil) != tt.wantErr {
				t.Errorf("RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("RemoveUser() there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_RemoveRefByUserIDLinkID(t *testing.T) {
	const expectedQuery = `DELETE FROM ref_link_user WHERE userID = (.+) AND linkid = (.+);`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "delete link",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(123, "linkid").
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: false,
		},
		{
			name: "delete link err",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WithArgs(123, "linkid").WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			db := sqlx.NewDb(baseDB, "postgres")
			defer db.Close()

			tt.prepare(mock)

			api := NewSQLAPI(db)
			if err := api.RemoveRefByUserIDLinkID(123, "linkid"); (err != nil) != tt.wantErr {
				t.Errorf("RemoveRefByUserIDLinkID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("RemoveRefByUserIDLinkID() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
