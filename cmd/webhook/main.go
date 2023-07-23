package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-academy-presentation/pkg/line"
	"go-academy-presentation/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
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
			line.Subscribe(userid)
		} else if e.Type == "unfollow" {
			line.Unsubscribe(userid)
		} else if e.Type == "message" {
			line.SendMessage(userid, text)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       response.ResponseBody("POST /webhook: OK"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
