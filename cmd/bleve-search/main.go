package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	_ "github.com/go-sql-driver/mysql"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	query, check := request.QueryStringParameters["query"]
	if !check {
		return events.APIGatewayProxyResponse{
			Body:       "GET /bleve-search: NG",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	languageCode := "ja"
	rows, err := db.Query("SELECT id, translated_description FROM garbage_item_details WHERE language_code = ?", languageCode)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	_, err = os.Stat("/tmp/ja.bluge")
	if err != nil {
		fmt.Println("4444444")
		panic(err.Error())
	}

	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()
	fieldMapping := bleve.NewTextFieldMapping()
	fieldMapping.Analyzer = cjk.AnalyzerName
	docMapping.AddFieldMappingsAt("Body", fieldMapping)
	indexMapping.AddDocumentMapping("Example", docMapping)

	// create a new index
	var index bleve.Index
	_, err = os.Stat("/tmp/ja.bleve")
	if err == nil {
		index, err = bleve.Open("/tmp/ja.bleve")
	} else {
		index, err = bleve.New("/tmp/ja.bleve", indexMapping)
	}
	if err != nil {
		fmt.Println("0000000")
		panic(err.Error())
	}

	// index some data
	data := make(map[string]string)
	for rows.Next() {
		var id string
		var description string
		if err := rows.Scan(&id, &description); err != nil {
			fmt.Println("1111111")
			panic(err.Error())
		}
		data[id] = description
	}
	for id, desc := range data {
		err := index.Index(id, desc)
		if err != nil {
			fmt.Println("3333333")
			panic(err.Error())
		}
	}

	q := bleve.NewMatchQuery(query)
	search := bleve.NewSearchRequest(q)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println("2222222")
		panic(err.Error())
	}

	for _, hit := range searchResults.Hits {
		fmt.Println(hit.ID)
		fmt.Println(hit.Score)
		doc := data[hit.ID]
		fmt.Println(doc)
	}

	return events.APIGatewayProxyResponse{
		Body:       "GET /search?query=" + query + ": OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
