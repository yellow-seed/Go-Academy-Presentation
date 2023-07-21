package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

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

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, garbage_id, classify FROM garbage_masters")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	stmt, err := db.Prepare("INSERT INTO garbage_items(garbage_id, garbage_master_id, category) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	for rows.Next() {
		var garbage_id string
		var garbage_master_id int
		var classify string

		err := rows.Scan(&garbage_master_id, &garbage_id, &classify)
		if err != nil {
			panic(err.Error())
		}

		var category string

		if classify == "可燃ごみ" {
			category = "burnable"
		} else if classify == "不燃ごみ" {
			category = "unburnable"
		} else if classify == "資源" {
			category = "recyclable"
		} else if classify == "粗大ごみ" {
			category = "large"
		} else {
			category = "other"
		}

		_, err = stmt.Exec(garbage_id, garbage_master_id, category)
		if err != nil {
			panic(err.Error())
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       "POST /garbageitem: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
