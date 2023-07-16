package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: 実体の実装
	q, check := request.QueryStringParameters["query"]

	if !check {
		return events.APIGatewayProxyResponse{
			Body:       "GET /bluge-search: NG",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "query=" + q + ": OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
