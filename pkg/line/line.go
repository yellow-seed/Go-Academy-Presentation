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

	if text == "switch language" {
		targetLang := ChangeLanguage(u.LanguageCode)
		u.LanguageCode = targetLang

		err = user.Update(db, u)
		if err != nil {
			panic(err.Error())
		}

		if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(message.SwitchLanguageMessage(targetLang))).Do(); err != nil {
			fmt.Println(err)
		}
	} else if text == "how to use" {
		if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(message.HowToUseMessage(u.LanguageCode))).Do(); err != nil {
			fmt.Println(err)
		}
	} else if text == "change search mode" {
		targetMode := ChangeSearchMode(u.SearchMode)
		u.SearchMode = targetMode

		err = user.Update(db, u)
		if err != nil {
			panic(err.Error())
		}

		if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(message.ChangeSearchModeMessage(u.LanguageCode, targetMode))).Do(); err != nil {
			fmt.Println(err)
		}
	} else {
		if u.SearchMode == "gpt" {
			if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(message.GPTSearchCautionMessage(u.LanguageCode))).Do(); err != nil {
				fmt.Println(err)
			}
		}
		res := search.Search(text, u.LanguageCode, u.SearchMode)
		if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(res)).Do(); err != nil {
			fmt.Println(err)
		}

		if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(message.ForeignerSupportMessage(u.LanguageCode))).Do(); err != nil {
			fmt.Println(err)
		}
	}
}

func Subscribe(userid string) {
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelToken"))
	if err != nil {
		panic(err.Error())
	}

	if _, err := bot.PushMessage(userid, linebot.NewTextMessage(message.HowToUseMessage("ja"))).Do(); err != nil {
		fmt.Println(err)
	}

	if _, err := bot.PushMessage(userid, linebot.NewTextMessage(message.InitialStateMessage())).Do(); err != nil {
		fmt.Println(err)
	}

	db, err := db.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = user.Read(db, userid)

	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからない場合は登録する
			u := &user.User{
				LineUserId:   userid,
				LanguageCode: "ja",
				SearchMode:   "sql",
			}
			err = user.Create(db, u)
			if err != nil {
				panic(err.Error())
			}

			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}
}

func Unsubscribe(userid string) {
	fmt.Println("unsubscribe")
	db, err := db.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	u, err := user.Read(db, userid)

	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからない場合なにもしない
		} else {
			panic(err.Error())
		}
	} else {
		err = user.Delete(db, u.LineUserId)
		if err != nil {
			panic(err.Error())
		}
	}
}

func ChangeLanguage(lang string) string {
	if lang == "en" {
		return "ja"
	} else if lang == "ja" {
		return "en"
	} else {
		return "ja"
	}
}

func ChangeSearchMode(mode string) string {
	if mode == "sql" {
		return "gpt"
	} else if mode == "gpt" {
		return "sql"
	} else {
		return "sql"
	}
}
