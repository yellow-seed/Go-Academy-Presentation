package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type GarbageItem struct {
	Id       int
	Category string
}

type GarbageMaster struct {
	GarbageId string
	Item      string
	Classify  string
	Remarks   string
}

type Garbage struct {
	GarbageMaster
	GarbageItem
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

	rows, err := db.Query("SELECT gi.id, gi.category, gm.garbage_id, gm.item, gm.classify, gm.remarks FROM garbage_masters gm INNER JOIN garbage_items gi ON gm.id = gi.garbage_master_id")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var garbages []Garbage
	for rows.Next() {
		var g Garbage
		if err := rows.Scan(&g.GarbageItem.Id, &g.GarbageItem.Category, &g.GarbageMaster.GarbageId, &g.GarbageMaster.Item, &g.GarbageMaster.Classify, &g.GarbageMaster.Remarks); err != nil {
			panic(err.Error())
		}
		garbages = append(garbages, g)
	}

	stmt, err := db.Prepare("INSERT INTO garbage_item_details(garbage_id, garbage_item_id, language_code, translated_name, translated_category, translated_description) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	for _, g := range garbages {
		var description string
		var translated_category string

		if g.GarbageItem.Category == "burnable" {
			translated_category = "可燃ごみ"
			description = g.GarbageMaster.Item + "は「" + translated_category + "」です。\n"
		} else if g.GarbageItem.Category == "unburnable" {
			translated_category = "不燃ごみ"
			description = g.GarbageMaster.Item + "は「" + translated_category + "」です。\n"
		} else if g.GarbageItem.Category == "recyclable" {
			translated_category = "資源"
			description = g.GarbageMaster.Item + "は「" + translated_category + "」です。\n"
		} else if g.GarbageItem.Category == "large" {
			translated_category = "粗大ごみ"
			description = g.GarbageMaster.Item + "は「" + translated_category + "」です。\n"
		} else {
			translated_category = "その他"
			description = g.GarbageMaster.Item + "は" + translated_category + "「" + g.GarbageMaster.Classify + "」です。\n"
		}

		if g.GarbageMaster.Remarks != "" {
			description += "[注意点]\n" + g.GarbageMaster.Remarks
		}

		_, err = stmt.Exec(g.GarbageMaster.GarbageId, g.GarbageItem.Id, "ja", g.GarbageMaster.Item, translated_category, description)
		if err != nil {
			panic(err.Error())
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "POST /garbageitemdetail-ja: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
