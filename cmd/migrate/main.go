package main

import (
	"go-academy-presentation/pkg/garbage"
	"go-academy-presentation/pkg/response"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	garbage.Migrate()

	return events.APIGatewayProxyResponse{
		Body:       response.ResponseBody("POST /migrate: OK"),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
