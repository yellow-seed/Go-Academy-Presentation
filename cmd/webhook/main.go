package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

	fmt.Println("-----")
	fmt.Println(event)
	fmt.Println("-----")

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
			postLineMessage(userid, text)
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

func postLineMessage(userid string, text string) {
	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelToken"))
	if err != nil {
		fmt.Println(err)
	}
	var res string

	// TODO: 特定のテキストの場合は処理を分岐する
	if text == "switch language" {
		res = text
	} else if text == "how to use" {
		res = text
	}

	if _, err := bot.PushMessage(userid, linebot.NewTextMessage(res)).Do(); err != nil {
		fmt.Println(err)
	}
}

func subscribe(userid string) {
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
			stmt, err := db.Prepare("INSERT INTO users(line_user_id, language) VALUES (?, ?)")
			if err != nil {
				panic(err.Error())
			}
			defer stmt.Close()

			_, err = stmt.Exec(userid, "ja")
			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	} else {
		// ユーザーがいる場合は何もしない
	}
}

func unsubscribe(userid string) {
	fmt.Println("unsubscribe")
}
