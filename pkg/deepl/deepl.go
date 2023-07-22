package deepl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func TranslateText(text, sourceLang, targetLang string) (string, error) {
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
