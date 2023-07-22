package user

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := &User{
			LineUserId:   "test",
			LanguageCode: "ja",
			SearchMode:   "sql",
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (line_user_id, language_code, search_mode) VALUES (?, ?, ?)")).WithArgs(u.LineUserId, u.LanguageCode, u.SearchMode).WillReturnResult(sqlmock.NewResult(1, 1))

		err = Create(db, u)
		if err != nil {
			t.Error(err.Error())
		}
	},
	)

	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		u := &User{
			LineUserId:   "test",
			LanguageCode: "ja",
			SearchMode:   "sql",
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (line_user_id, language_code, search_mode) VALUES (?, ?, ?)")).WithArgs(u.LineUserId, u.LanguageCode, u.SearchMode).WillReturnResult(sqlmock.NewErrorResult(errors.New("ERROR!!!"))).WillReturnError(errors.New("INSERT FAILED!!!"))

		err = Create(db, u)
		if err == nil {
			t.Error(err.Error())
		}
	},
	)
}

func TestRead(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		lineUserId := "test"

		rows := sqlmock.NewRows([]string{"id", "line_user_id", "language_code", "search_mode"}).AddRow(1, lineUserId, "ja", "sql")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?")).WithArgs(lineUserId).WillReturnRows(rows)

		_, err = Read(db, lineUserId)
		if err != nil {
			t.Error(err.Error())
		}
	},
	)

	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		lineUserId := "test"

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?")).WithArgs(lineUserId).WillReturnError(errors.New("ERROR!!!"))

		_, err = Read(db, lineUserId)
		if err == nil {
			t.Error(err.Error())
		}
	},
	)
}

func TestUpdate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		u := &User{
			Id:           1,
			LineUserId:   "test",
			LanguageCode: "ja",
			SearchMode:   "sql",
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET language_code = ?, search_mode = ? WHERE id = ?")).WithArgs(u.LanguageCode, u.SearchMode, u.Id).WillReturnResult(sqlmock.NewResult(1, 1))

		err = Update(db, u)
		if err != nil {
			t.Error(err.Error())
		}
	},
	)

	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		u := &User{
			Id:           1,
			LineUserId:   "test",
			LanguageCode: "ja",
			SearchMode:   "sql",
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET language_code = ?, search_mode = ? WHERE id = ?")).WithArgs(u.LanguageCode, u.SearchMode, u.Id).WillReturnError(errors.New("ERROR!!!"))

		err = Update(db, u)
		if err == nil {
			t.Error(err.Error())
		}
	},
	)
}

func TestDelete(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		lineUserId := "test"

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE line_user_id = ?")).WithArgs(lineUserId).WillReturnResult(sqlmock.NewResult(1, 1))

		err = Delete(db, lineUserId)
		if err != nil {
			t.Error(err.Error())
		}
	},
	)

	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		lineUserId := "test"

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE line_user_id = ?")).WithArgs(lineUserId).WillReturnError(errors.New("ERROR!!!"))

		err = Delete(db, lineUserId)
		if err == nil {
			t.Error(err.Error())
		}
	},
	)
}
