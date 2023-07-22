package user

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id           int
	LineUserId   string
	LanguageCode string
	SearchMode   string
}

func Create(db *sql.DB, u *User) error {
	_, err := db.Exec("INSERT INTO users (line_user_id, language_code, search_mode) VALUES (?, ?, ?)", u.LineUserId, u.LanguageCode, u.SearchMode)
	if err != nil {
		return err
	}
	return nil
}

func Read(db *sql.DB, lineUserId string) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?", lineUserId).Scan(&u.Id, &u.LineUserId, &u.LanguageCode, &u.SearchMode)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func Update(db *sql.DB, u *User) error {
	_, err := db.Exec("UPDATE users SET language_code = ?, search_mode = ? WHERE id = ?", u.LanguageCode, u.SearchMode, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func Delete(db *sql.DB, lineUserId string) error {
	_, err := db.Exec("DELETE FROM users WHERE line_user_id = ?", lineUserId)
	if err != nil {
		return err
	}
	return nil
}
