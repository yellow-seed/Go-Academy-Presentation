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
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")
	lineId := os.Getenv("LineUserId")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO users(line_user_id, language_code, search_mode) VALUES (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	insert.Exec(lineId, "ja", "sql")

	return events.APIGatewayProxyResponse{
		Body:       "POST /masteruser: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
