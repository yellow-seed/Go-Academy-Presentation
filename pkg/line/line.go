package line

import (
	"database/sql"
	"fmt"
	"go-academy-presentation/pkg/db"
	"go-academy-presentation/pkg/message"
	"go-academy-presentation/pkg/search"
	"go-academy-presentation/pkg/user"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func SendMessage(userid string, text string) {
	db, err := db.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	u, err := user.Read(db, userid)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("-----")
	fmt.Println(u.Id)
	fmt.Println(u.LineUserId)
	fmt.Println(u.LanguageCode)
	fmt.Println(u.SearchMode)
	fmt.Println("-----")

	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelToken"))
	if err != nil {
		panic(err.Error())
	}

	var res string
	if text == "switch language" {
		targetLang := changeLanguage(u.LanguageCode)
		u.LanguageCode = targetLang

		err = user.Update(db, u)
		if err != nil {
			panic(err.Error())
		}

		res = message.SwitchLanguageMessage(targetLang)
	} else if text == "how to use" {
		res = message.HowToUseMessage(u.LanguageCode)
	} else if text == "change search mode" {
		targetMode := changeSearchMode(u.SearchMode)
		u.SearchMode = targetMode

		err = user.Update(db, u)
		if err != nil {
			panic(err.Error())
		}

		res = message.ChangeSearchModeMessage(u.LanguageCode, targetMode)
	} else {
		// TODO: GPT検索の場合、実行前に一度メッセージを送る
		res = search.Search(text, u.LanguageCode, u.SearchMode)
	}

	if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(res)).Do(); err != nil {
		fmt.Println(err)
	}

	// TODO: 日本語以外の場合は追加のメッセージを送る　外国人向けのページのリンク付きのメッセージ
}

// func SendMessage(userid string, text string) {
// 	// var u user.User
// 	// dbuser := os.Getenv("DBUser")
// 	// pass := os.Getenv("DBPass")
// 	// host := os.Getenv("DBHost")
// 	// name := os.Getenv("DBName")
// 	// db, err := sql.Open("mysql", dbuser+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	// defer db.Close()

// 	db, err := db.InitDB()
// 	if err != nil {
// 		fmt.Println("------")
// 		fmt.Println(err)
// 		panic(err.Error())
// 	}
// 	defer db.Close()

// 	// if err = db.QueryRow("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?", userid).Scan(&u.Id, &u.LineUserId, &u.LanguageCode, &u.SearchMode); err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	u, err := user.Read(db, userid)
// 	if err != nil {
// 		fmt.Println("======")
// 		fmt.Println(err)
// 		panic(err.Error())
// 	}

// 	fmt.Println("-----")
// 	fmt.Println(u.Id)
// 	fmt.Println(u.LineUserId)
// 	fmt.Println(u.LanguageCode)
// 	fmt.Println(u.SearchMode)
// 	fmt.Println("-----")

// 	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelToken"))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	var res string
// 	if text == "switch language" {
// 		targetLang := changeLanguage(u.LanguageCode)
// 		u.LanguageCode = targetLang
// 		err := user.Update(db, u)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		// update, err := db.Prepare("UPDATE users SET language_code = ? WHERE id = ?")
// 		// if err != nil {
// 		// 	panic(err.Error())
// 		// }
// 		// update.Exec(targetLang, u.Id)
// 		res = message.SwitchLanguageMessage(targetLang)
// 	} else if text == "how to use" {
// 		res = message.HowToUseMessage(u.LanguageCode)
// 	} else if text == "change search mode" {
// 		update, err := db.Prepare("UPDATE users SET search_mode = ? WHERE id = ?")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		targetMode := changeSearchMode(u.SearchMode)
// 		update.Exec(targetMode, u.Id)
// 		res = message.ChangeSearchModeMessage(u.LanguageCode, targetMode)
// 	} else {
// 		// TODO: GPT検索の場合、実行前に一度メッセージを送る
// 		res = search.Search(text, u.LanguageCode, u.SearchMode)
// 	}

// 	if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(res)).Do(); err != nil {
// 		fmt.Println(err)
// 	}

// 	// TODO: 日本語以外の場合は追加のメッセージを送る　外国人向けのページのリンク付きのメッセージ
// }

func Subscribe(userid string) {
	// TODO: 登録時専用のメッセージを投稿する
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.QueryRow("SELECT line_user_id FROM users WHERE line_user_id = ?", userid).Scan(&userid)

	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからない場合は登録する
			stmt, err := db.Prepare("INSERT INTO users(line_user_id, language_code, search_mode) VALUES (?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}
			defer stmt.Close()

			_, err = stmt.Exec(userid, "ja", "sql")
			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}
}

func Unsubscribe(userid string) {
	// TODO: ユーザー削除処理
	fmt.Println("unsubscribe")
}

func changeLanguage(lang string) string {
	if lang == "en" {
		return "ja"
	} else if lang == "ja" {
		return "en"
	} else {
		return "ja"
	}
}

func changeSearchMode(mode string) string {
	if mode == "sql" {
		return "gpt"
	} else if mode == "gpt" {
		return "sql"
	} else {
		return "sql"
	}
}
