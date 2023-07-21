package user

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id           int    
	LineUserId   string 
	LanguageCode string 
	SearchMode   string 
}

func Create(db *sql.DB, u *User) error {
	insert, err := db.Prepare("INSERT INTO users (line_user_id, language_code, search_mode) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	insert.Exec(u.LineUserId, u.LanguageCode, u.SearchMode)
	return nil
}

func Read(db *sql.DB, lineUserId string) (*User, error) {
	u := &User{}
	err := db.QueryRow("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?", lineUserId).Scan(&u.Id, &u.LineUserId, &u.LanguageCode, &u.SearchMode)
	if err != nil {
		fmt.Println("||||||||||")
		fmt.Println(err)
		return nil, err
	}
	return u, nil
}

func Update(db *sql.DB, u *User) error {
	update, err := db.Prepare("UPDATE users SET language_code = ?, search_mode = ? WHERE id = ?")
	if err != nil {
		return err
	}
	update.Exec(u.LanguageCode, u.SearchMode, u.Id)
	return nil
}

func Delete(db *sql.DB, lineUserId string) error {
	delete, err := db.Prepare("DELETE FROM users WHERE line_user_id = ?")
	if err != nil {
		return err
	}
	delete.Exec(lineUserId)
	return nil
}
