package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type GarbageItemDetail struct {
	Id                    int
	GarbageId             string
	GarbageItemId         int
	LanguageCode          string
	TranslatedName        string
	TranslatedCategory    string
	TranslatedDescription string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	q, queryCheck := request.QueryStringParameters["query"]
	l, langCheck := request.QueryStringParameters["lang"]

	if !queryCheck {
		return events.APIGatewayProxyResponse{
			Body:       "GET /bluge-search: NG",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if !langCheck {
		l = "ja"
	}

	// TODO: GPT対応ができるようになったら切り替える
	result := likeSearch(q, l)

	return events.APIGatewayProxyResponse{
		Body:       "result=\n" + result + "\n: OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func likeSearch(q string, lang string) string {
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	q = "%" + q + "%"
	rows, err := db.Query("SELECT id, garbage_id, garbage_item_id, language_code, translated_name, translated_category, translated_description FROM garbage_item_details WHERE language_code = ? AND translated_description LIKE ?", lang, q)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var details []GarbageItemDetail
	for rows.Next() {
		var detail GarbageItemDetail
		if err := rows.Scan(&detail.Id, &detail.GarbageId, &detail.GarbageItemId, &detail.LanguageCode, &detail.TranslatedName, &detail.TranslatedCategory, &detail.TranslatedDescription); err != nil {
			panic(err.Error())
		}
		details = append(details, detail)
	}
	var t string

	if len(details) > 24 {
		if lang == "en" {
			t = "There are more than 25 search results. Show some.\n"
		} else {
			t = "25件以上の検索結果があります。一部を表示します。\n"
		}

		for i := 0; i < 25; i++ {
			t += details[i].TranslatedName + "\n"
		}
	} else if len(details) >= 10 {
		if lang == "en" {
			t = "There are " + strconv.Itoa(len(details)) + " search results.\n"
		} else {
			t = strconv.Itoa(len(details)) + "件の検索結果があります。\n"
		}
		for i := 0; i < len(details); i++ {
			t += details[i].TranslatedName + "\n"
		}
	} else {
		if lang == "en" {
			t = "There are " + strconv.Itoa(len(details)) + " search results.\n"
		} else {
			t = strconv.Itoa(len(details)) + "件の検索結果があります。\n"
		}
		for i := 0; i < len(details); i++ {
			t += details[i].TranslatedDescription + "\n"
		}
	}
	return t
}
