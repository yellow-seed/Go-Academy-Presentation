package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         int
	LineUserId string
	Language   string
}

type GarbageMaster struct {
	Id         int
	PublicCode string
	GarbageId  string
	PublicName string
	District   string
	Item       string
	ItemKana   string
	ItemEng    string
	Classify   string
	Note       string
	Remarks    string
	LargeFee   string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	lineId := os.Getenv("LINE_USER_ID")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO users(line_user_id, language) VALUES (?, ?)")

	if err != nil {
		panic(err.Error())
	}

	insert.Exec(lineId, "en")

	return events.APIGatewayProxyResponse{
		Body:       "POST /masteruser: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
