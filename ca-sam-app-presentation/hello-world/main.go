package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

// リクエストボディを受け取る構造体
type Response struct {
	RequestBody string `json:"RequestBody"`
}

// リクエストボディから特定のパラメータを受け取る構造体
type Event struct {
	Events []struct {
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
	// TODO: type followのときはDBにユーザーを登録する
	// https://developers.line.biz/ja/docs/messaging-api/getting-user-ids/#get-user-ids-in-webhook

	// 先ほど定義した構造体
	var event Event

	// リクエストボディを受け取る
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
			Body: "Hello"}, nil
	}

	// 送信元userIDとテキスト内容を変数に代入
	userid := fmt.Sprintf("%v", event.Events[0].Source.UserId)
	text := fmt.Sprintf("%v", event.Events[0].Message.Text)

	// 送信元ユーザにメッセージを返信する関数(次に実装します)
	postLineMessage(userid, text)

	return events.APIGatewayProxyResponse{
		Body:       "Hello",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func postLineMessage(userid string, text string) {
	err := godotenv.Load(".env")

	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}

	if _, err := bot.PushMessage(userid, linebot.NewTextMessage(text)).Do(); err != nil {
		fmt.Println(err)
	}
}
