package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/blugelabs/bluge"
	segment "github.com/blugelabs/bluge_segment_api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ikawaha/blugeplugin/analysis/lang/ja"
)

type Document struct {
	ID string
	Text string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	q, check := request.QueryStringParameters["query"]
	if !check {
		return events.APIGatewayProxyResponse{
			Body:       "GET /bluge-search: NG",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	fmt.Print("---------")
	config := bluge.DefaultConfig("/tmp/ja.bluge")
	w, err := bluge.OpenWriter(config)
	if err != nil {
		fmt.Println("0000000")
		log.Fatalf("error opening writer: %v", err)
	}
	defer w.Close()

	fmt.Println("==========")
	docs := NewDocuments()

	fmt.Println("++++++++++")
	//indexing
	for _, doc := range docs {
		doc := doc
		if err := w.Update(doc.ID(), doc); err != nil {
			fmt.Println("1111111")
			log.Fatalf("error updating document: %v", err)
		}

		// display
		fmt.Printf("indexed document with id:%s\n", doc.ID())
		doc.EachField(func(field segment.Field) {
			fmt.Printf("\t%s: %s\n", field.Name(), field.Value())
		})
	}

	fmt.Println("**********")
	// index reader
	r, err := w.Reader()
	if err != nil {
		fmt.Println("2222222")
		log.Fatalf("error getting index reader: %v", err)
	}
	defer r.Close()

	query := bluge.NewMatchQuery(q).SetAnalyzer(ja.Analyzer()).SetField("body")
	req := bluge.NewTopNSearch(10, query).WithStandardAggregations()
	fmt.Printf("query: search field %q, value %q\n", query.Field(), q)

	// search
	ite, err := r.Search(context.Background(), req)
	if err != nil {
		fmt.Println("3333333")
		log.Fatalf("error executing search: %v", err)
	}

	// search result
	for {
		match, err := ite.Next()
		if err != nil {
			fmt.Println("4444444")
			log.Fatalf("error iterating search results: %v", err)
		}
		if match == nil {
			break
		}
		if err := match.VisitStoredFields(func(field string, value []byte) bool{
			fmt.Printf("%s: %q\n", field, string(value))
			return true
		}); err != nil {
			fmt.Println("5555555")
			log.Fatalf("error visiting stored fields: %v", err)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "GET /search?query=" + q + ": OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func NewDocuments() []*bluge.Document {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		fmt.Println("6666666")
		panic(err.Error())
	}
	defer db.Close()

	languageCode := "ja"
	rows, err := db.Query("SELECT id, translated_description FROM garbage_item_details WHERE language_code = ?", languageCode)
	if err != nil {
		fmt.Println("7777777")
		panic(err.Error())
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.Text)
		if err != nil {
			fmt.Println("8888888")
			panic(err.Error())
		}

		documents = append(documents, doc)
	}

	var ret []*bluge.Document
	for _, v := range documents {
		id := bluge.NewTextField("id", v.ID).WithAnalyzer(ja.Analyzer())
		body := bluge.NewTextField("body", v.Text).WithAnalyzer(ja.Analyzer())
		doc := bluge.NewDocument(v.ID).AddField(id).AddField(body)
		ret = append(ret, doc)
	}
	return ret
}
