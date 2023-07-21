package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
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

type Data struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Result struct {
	IsExist bool   `json:"result"`
	Data    []Data `json:"data"`
}

func Search(q string, lang string, mode string) string {
	var res string
	fmt.Println("Search")
	fmt.Println("-----")
	fmt.Println(q)
	fmt.Println(lang)
	fmt.Println(mode)
	fmt.Println("-----")
	if mode == "sql" {
		res = SqlLikeSearch(q, lang)
	} else if mode == "gpt" {
		res = GptSearch(q, lang)
	} else {
		res = SqlLikeSearch(q, lang)
	}
	return res
}

func SqlLikeSearch(q string, lang string) string {
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	fmt.Println("SqlLikeSearch")
	fmt.Println("-----")
	fmt.Println(user)
	fmt.Println(pass)
	fmt.Println(host)
	fmt.Println(name)
	fmt.Println("-----")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "%" + q + "%"
	rows, err := db.Query("SELECT id, garbage_id, garbage_item_id, language_code, translated_name, translated_category, translated_description FROM garbage_item_details WHERE language_code = ? AND translated_description LIKE ?", lang, query)
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

func GptSearch(q string, lang string) string {
	// TODO: 対象データ数をもう少し絞ってもいいかも
	fmt.Println("GPT Search")
	data := CreateData(lang)
	fmt.Println("CreateData")
	limit := 20
	partSize := len(data) / limit

	var resultData []Data

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * partSize
			end := start + partSize
			if i == limit-1 {
				end = len(data)
			}
			partData := data[start:end]
			jsonData, err := json.Marshal(partData)
			if err != nil {
				fmt.Println(err)
			}
			var r Result
			fmt.Println(string(jsonData))
			// t := LangchainNameSearch(q, string(jsonData))
			t := LangchainDescriptionSearch(q, string(jsonData))
			err = json.Unmarshal([]byte(t), &r)

			if err == nil && r.IsExist {
				mutex.Lock()
				resultData = append(resultData, r.Data...)
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	var t string

	for _, d := range resultData {
		fmt.Println("=====")
		fmt.Println(d.Id)
		fmt.Println(d.Description)
	}

	if len(resultData) > 24 {
		if lang == "en" {
			t = "There are more than 25 search results. Show some.\n"
		} else {
			t = "25件以上の検索結果があります。一部を表示します。\n"
		}

		for i := 0; i < 25; i++ {
			t += resultData[i].Name + "\n"
		}
	} else if len(resultData) >= 10 {
		if lang == "en" {
			t = "There are " + strconv.Itoa(len(resultData)) + " search results.\n"
		} else {
			t = strconv.Itoa(len(resultData)) + "件の検索結果があります。\n"
		}
		for i := 0; i < len(resultData); i++ {
			t += resultData[i].Description + "\n"
		}
	} else {
		if lang == "en" {
			t = "There are " + strconv.Itoa(len(resultData)) + " search results.\n"
		} else {
			t = strconv.Itoa(len(resultData)) + "件の検索結果があります。\n"
		}
		for i := 0; i < len(resultData); i++ {
			t += resultData[i].Description + "\n"
		}
	}
	return t
}

func CreateData(lang string) []Data {
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, translated_name, translated_description FROM garbage_item_details WHERE language_code = ?", lang)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var data []Data
	for rows.Next() {
		var d Data
		if err := rows.Scan(&d.Id, &d.Name, &d.Description); err != nil {
			fmt.Println(err.Error())
		}
		// log.Fatalln(d.Name)
		data = append(data, d)
	}
	return data
}

func LangchainNameSearch(keyword string, stringData string) string {
	baseCommand :=
		`
キーワード「%s」で検索を行う。
次のフォーマットで値を抽出せよ。

## Output Format

{
  "result": $[意味の近いnameの値があればtrue, なければfalseを返してください]
  "data": [
    {
      "id": $[意味の近いデータのID 1番目],
			"name": $[意味の近いデータの名前 1番目],
      "description": $[意味の近いデータの説明文 1番目]
    },
    ...
    {
      "id": $[意味の近いデータのID N番目],
			"name": $[意味の近いデータの名前 N番目],
      "description": $[意味の近いデータの説明文 N番目]
    },
  ]
}
キーは必ず含ませる。
JSON以外の情報は削除する。
id, name, descriptionは必ず元のテキストに含まれる文字列だけを値として使う。
該当する情報がない場合 null にする。
`
	command := fmt.Sprintf(baseCommand, keyword)

	emptyResponse :=
		`
{
	"result": false,
	"data": []
}
`
	baseAsk := `
"data": %s 
`
	ask := fmt.Sprintf(baseAsk, stringData)

	answer :=
		`
{
	"result": $[意味の近いnameの値があればtrue, なければfalseを返してください]
	"comment": $[JSON形式ではない内容はここに記入してください]
	"data": [
		{
			"id": $[意味の近いデータのID 1番目],
			"name": $[意味の近いデータの名前 1番目],
			"description": $[意味の近いデータの説明文 1番目]
		},
		...
		{
			"id": $[意味の近いデータのID N番目],
			"name": $[意味の近いデータの名前 N番目],
			"description": $[意味の近いデータの説明文 N番目]
		},
	]
}
`
	systemMessage :=
		`
{
	"role": "system",
	"content": %s,
},
{
	"role": "user",
	"content": "",
},
{
	"role": "assistant",
	"content": %s,
},
{
	"role": "user",
	"content": %s,
},
{
	"role": "assistant",
	"content": %s,
}
`
	s := fmt.Sprintf(systemMessage, command, emptyResponse, ask, answer)

	llm, err := openai.NewChat(openai.WithToken(os.Getenv("OPENAIApiKey")))
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Text: s},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}),
		llms.WithTemperature(0.0),
		llms.WithModel("gpt-3.5-turbo-0613"),
	)
	if err != nil {
		fmt.Println(err)
	}

	return completion
}

func LangchainDescriptionSearch(keyword string, stringData string) string {
	baseCommand :=
		`
キーワード「%s」で検索を行う。
次のフォーマットで値を抽出せよ。

## Output Format

{
  "result": $[意味の近いdescriptionの値があればtrue, なければfalseを返してください]
  "data": [
    {
      "id": $[意味の近いデータのID 1番目],
			"name": $[意味の近いデータの名前 1番目],
      "description": $[意味の近いデータの説明文 1番目]
    },
    ...
    {
      "id": $[意味の近いデータのID N番目],
			"name": $[意味の近いデータの名前 N番目],
      "description": $[意味の近いデータの説明文 N番目]
    },
  ]
}
キーは必ず含ませる。
JSON以外の情報は削除する。
id, name, descriptionは必ず元のテキストに含まれる文字列だけを値として使う。
該当する情報がない場合 null にする。
`
	command := fmt.Sprintf(baseCommand, keyword)

	emptyResponse :=
		`
{
	"result": false,
	"data": []
}
`
	baseAsk := `
"data": %s 
`
	ask := fmt.Sprintf(baseAsk, stringData)

	answer :=
		`
{
	"result": $[意味の近いdescriptionの値があればtrue, なければfalseを返してください]
	"comment": $[JSON形式ではない内容はここに記入してください]
	"data": [
		{
			"id": $[意味の近いデータのID 1番目],
			"name": $[意味の近いデータの名前 1番目],
			"description": $[意味の近いデータの説明文 1番目]
		},
		...
		{
			"id": $[意味の近いデータのID N番目],
			"name": $[意味の近いデータの名前 N番目],
			"description": $[意味の近いデータの説明文 N番目]
		},
	]
}
`
	systemMessage :=
		`
{
	"role": "system",
	"content": %s,
},
{
	"role": "user",
	"content": "",
},
{
	"role": "assistant",
	"content": %s,
},
{
	"role": "user",
	"content": %s,
},
{
	"role": "assistant",
	"content": %s,
}
`
	s := fmt.Sprintf(systemMessage, command, emptyResponse, ask, answer)

	llm, err := openai.NewChat(openai.WithToken(os.Getenv("OPENAIApiKey")))
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Text: s},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}),
		llms.WithTemperature(0.0),
		llms.WithModel("gpt-3.5-turbo-0613"),
		// llms.WithModel("gpt-4-0613"),
	)
	if err != nil {
		fmt.Println(err)
	}

	return completion
}
