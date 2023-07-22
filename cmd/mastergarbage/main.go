package main

import (
	"encoding/csv"
	"go-academy-presentation/pkg/db"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

const (
	url = "https://www.opendata.metro.tokyo.lg.jp/setagaya/131121_setagayaku_garbage_separate.csv"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := db.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

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
	os.Remove(path.Base(url))

	return events.APIGatewayProxyResponse{
		Body:       "POST /mastergarbage: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
