package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-academy-presentation/pkg/message"
	"go-academy-presentation/pkg/search"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
)

// リクエストボディを受け取る構造体
type Response struct {
	RequestBody string `json:"RequestBody"`
}

// リクエストボディから特定のパラメータを受け取る構造体
type Event struct {
	Events []struct {
		Type    string `json:"type"`
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
		Source struct {
			UserId string `json:"userId"`
		} `json:"source"`
	}
}

type User struct {
	Id           int    `json:"id"`
	LineUserId   string `json:"line_user_id"`
	LanguageCode string `json:"language_code"`
	SearchMode   string `json:"search_mode"`
}

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

// func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	resp, err := http.Get(DefaultHTTPGetAddress)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{}, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return events.APIGatewayProxyResponse{}, ErrNon200Response
// 	}

// 	ip, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{}, err
// 	}

// 	if len(ip) == 0 {
// 		return events.APIGatewayProxyResponse{}, ErrNoIP
// 	}

// 	return events.APIGatewayProxyResponse{
// 		Body:       fmt.Sprintf("Hello, %v", string(ip)),
// 		StatusCode: 200,
// 	}, nil
// }

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event Event

	body := request.Body
	res := Response{
		RequestBody: body,
	}

	json.Unmarshal([]byte(res.RequestBody), &event)

	if len(event.Events) == 0 {
		fmt.Println("event is empty")
		return events.APIGatewayProxyResponse{
			Body:       "Hello",
			StatusCode: 200,
		}, nil
	}

	for _, e := range event.Events {
		var userid string
		if e.Source.UserId != "" {
			userid = fmt.Sprintf("%v", e.Source.UserId)
		} else {
			continue
		}
		text := fmt.Sprintf("%v", e.Message.Text)

		if e.Type == "follow" {
			subscribe(userid)
		} else if e.Type == "unfollow" {
			unsubscribe(userid)
		} else if e.Type == "message" {
			var u User
			user := os.Getenv("DBUser")
			pass := os.Getenv("DBPass")
			host := os.Getenv("DBHost")
			name := os.Getenv("DBName")

			db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()

			if err = db.QueryRow("SELECT id, line_user_id, language_code, search_mode FROM users WHERE line_user_id = ?", userid).Scan(&u.Id, &u.LineUserId, &u.LanguageCode, &u.SearchMode); err != nil {
				fmt.Println(err)
			}

			fmt.Println("-----")
			fmt.Println(u.Id)
			fmt.Println(u.LineUserId)
			fmt.Println(u.LanguageCode)
			fmt.Println(u.SearchMode)
			fmt.Println("-----")

			postLineMessage(u, text)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "POST /webhook: OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func postLineMessage(u User, text string) {
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelToken"))
	if err != nil {
		fmt.Println(err)
	}
	var res string

	if text == "switch language" {
		user := os.Getenv("DBUser")
		pass := os.Getenv("DBPass")
		host := os.Getenv("DBHost")
		name := os.Getenv("DBName")

		db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		update, err := db.Prepare("UPDATE users SET language_code = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		targetLang := changeLanguage(u.LanguageCode)
		update.Exec(targetLang, u.Id)
		res = message.SwitchLanguageMessage(targetLang)
	} else if text == "how to use" {
		res = message.HowToUseMessage(u.LanguageCode)
	} else if text == "change search mode" {
		user := os.Getenv("DBUser")
		pass := os.Getenv("DBPass")
		host := os.Getenv("DBHost")
		name := os.Getenv("DBName")

		db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		update, err := db.Prepare("UPDATE users SET search_mode = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		targetMode := changeSearchMode(u.SearchMode)
		update.Exec(targetMode, u.Id)
		res = message.ChangeSearchModeMessage(u.LanguageCode, targetMode)
	} else {
		// TODO: GPT検索の場合、実行前に一度メッセージを送る
		res = search.Search(text, u.LanguageCode, u.SearchMode)
		// if u.SearchMode == "sql" {
		// 	res = search.SqlLikeSearch(text, u.LanguageCode)
		// } else {
		// 	res = search.GptSearch(text, u.LanguageCode)
		// }
	}

	if _, err := bot.PushMessage(u.LineUserId, linebot.NewTextMessage(res)).Do(); err != nil {
		fmt.Println(err)
	}

	// TODO: 日本語以外の場合は追加のメッセージを送る　外国人向けのページのリンク付きのメッセージ
}

func subscribe(userid string) {
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
			// ユーザーがいない場合は登録する
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

func unsubscribe(userid string) {
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
