package main

import (
	"go-academy-presentation/pkg/search"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	q, queryCheck := request.QueryStringParameters["query"]
	l, langCheck := request.QueryStringParameters["lang"]
	m, modeCheck := request.QueryStringParameters["mode"]

	if !queryCheck {
		return events.APIGatewayProxyResponse{
			Body:       "GET /search: NG",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if !langCheck {
		l = "ja"
	}

	if !modeCheck {
		m = "sql"
	}

	result := search.Search(q, l, m)

	return events.APIGatewayProxyResponse{
		Body:       "result=\n" + result + "\n: OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
