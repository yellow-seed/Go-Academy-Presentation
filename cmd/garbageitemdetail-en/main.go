package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	ItemEng   string
	Classify  string
	Remarks   string
}

type GarbageItemDetail struct {
	GarbageId             string
	GarbageItemId         int
	LanguageCode          string
	TranslatedName        string
	TranslatedCategory    string
	TranslatedDescription string
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

	rows, err := db.Query("SELECT gi.id, gi.category, gm.garbage_id, gm.item_eng, gm.classify, gm.remarks FROM garbage_masters gm INNER JOIN garbage_items gi ON gm.id = gi.garbage_master_id")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var garbages []Garbage
	var details []GarbageItemDetail
	for rows.Next() {
		var g Garbage
		if err := rows.Scan(&g.GarbageItem.Id, &g.GarbageItem.Category, &g.GarbageMaster.GarbageId, &g.GarbageMaster.ItemEng, &g.GarbageMaster.Classify, &g.GarbageMaster.Remarks); err != nil {
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
		var detail GarbageItemDetail
		var translated_description string
		var translated_category string
		var translated_remarks string

		if g.GarbageMaster.Remarks != "" {
			translated_remarks, err = translateText(g.GarbageMaster.Remarks, "ja", "en")
			if err != nil {
				panic(err.Error())
			}
		}

		if g.GarbageItem.Category == "burnable" {
			translated_category = "burnable garbage"
			translated_description = g.GarbageMaster.ItemEng + `is "` + translated_category + `".\n`
		} else if g.GarbageItem.Category == "unburnable" {
			translated_category = "non-burnable garbage"
			translated_description = g.GarbageMaster.ItemEng + `is "` + translated_category + `".\n`
		} else if g.GarbageItem.Category == "recyclable" {
			translated_category = "recyclable waste"
			translated_description = g.GarbageMaster.ItemEng + `is "` + translated_category + `".\n`
		} else if g.GarbageItem.Category == "large" {
			translated_category = "oversized garbage"
			translated_description = g.GarbageMaster.ItemEng + `is "` + translated_category + `".\n`
		} else {
			translated_category = "other"
			other_classify, err := translateText(g.GarbageMaster.Classify, "ja", "en")
			if err != nil {
				fmt.Println(err)
			}
			translated_description = g.GarbageMaster.ItemEng + `is` + translated_category + `"` + other_classify + `".\n`
		}

		if g.GarbageMaster.Remarks != "" {
			translated_description += "[Note]\n" + translated_remarks
		}

		detail.GarbageId = g.GarbageMaster.GarbageId
		detail.GarbageItemId = g.GarbageItem.Id
		detail.LanguageCode = "en"
		detail.TranslatedName = g.GarbageMaster.ItemEng
		detail.TranslatedCategory = translated_category
		detail.TranslatedDescription = translated_description

		details = append(details, detail)
	}

	for _, d := range details {
		_, err = stmt.Exec(d.GarbageId, d.GarbageItemId, d.LanguageCode, d.TranslatedName, d.TranslatedCategory, d.TranslatedDescription)
		if err != nil {
			panic(err.Error())
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "POST /garbageitemdetail-en: OK",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func translateText(text, sourceLang, targetLang string) (string, error) {
	// 翻訳APIのURLとパラメータ
	apiURL := "https://api.deepl.com/v2/translate"
	authKey := os.Getenv("DEEPLApiKey")
	data := url.Values{}
	data.Set("text", text)
	data.Set("target_lang", targetLang)

	// 翻訳APIにリクエストを送信
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to send translation request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", authKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to send translation request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのJSONをパースして翻訳結果を取得
	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse translation response: %v", err)
	}
	if len(result.Translations) == 0 {
		return "", fmt.Errorf("no translations found")
	}
	return result.Translations[0].Text, nil
}
