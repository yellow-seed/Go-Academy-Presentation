package main

import (
	"database/sql"
	"encoding/csv"
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

	url := "https://www.opendata.metro.tokyo.lg.jp/setagaya/131121_setagayaku_garbage_separate.csv"

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	r := csv.NewReader(res.Body)

	// headerを読み飛ばす
	_, err = r.Read()
	if err != nil {
		panic(err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	stmt, err := tx.Prepare("INSERT INTO garbage_masters(public_code, garbage_id, public_name, district, item, item_kana, item_eng, classify, note, remarks, large_fee) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		if err != nil {
			panic(err.Error())
		}
	}

	tx.Commit()

	// delete download file
	os.Remove("131121_setagayaku_garbage_separate.csv")

	return events.APIGatewayProxyResponse{
		Body:       "POST /mastergarbage: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
