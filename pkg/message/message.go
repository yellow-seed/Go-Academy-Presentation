package message

func SwitchLanguageMessage(targetLang string) string {
	var res string
	if targetLang == "ja" {
		res = "言語を日本語に変更しました。"
	} else if targetLang == "en" {
		res = "Changed the language to English."
	}
	return res
}

func HowToUseMessage(lang string) string {
	var res string
	if lang == "ja" {
		res = "世田谷区のごみの分類方法を提供するbotです。分別方法を知りたいごみの名前を送信してください。文字一致検索モードとChatGPTを用いた類似単語検索モードがあります。検索モードの切り替えはメニューのボタンからおこなえます。"
	} else if lang == "en" {
		res = "This is a bot that provides a method of sorting garbage in Setagaya Ward. Please send us the name of the garbage you want to know how to separate. There are character matching search mode and similar word search mode using ChatGPT. You can switch the search mode from the button on the menu."
	}
	return res
}

func ChangeSearchModeMessage(lang string, targetSearchMode string) string {
	var res string
	if lang == "ja" {
		res = "検索モードを" + SearchModeName(targetSearchMode, lang) + "に切り替えました。"
	} else if lang == "en" {
		res = "Changed the search mode to " + SearchModeName(targetSearchMode, lang) + "."
	} else {
		res = "検索モードを" + SearchModeName(targetSearchMode, lang) + "に切り替えました。"
	}
	return res
}

func SearchModeName(mode string, lang string) string {
	if mode == "sql" {
		return SQLSearchModeName(lang)
	} else if mode == "gpt" {
		return GPTSearchModeName(lang)
	} else {
		return SQLSearchModeName(lang)
	}
}

func SQLSearchModeName(lang string) string {
	var res string
	if lang == "ja" {
		res = "文字一致検索"
	} else if lang == "en" {
		res = "Exact match search"
	} else {
		res = "文字一致検索"
	}
	return res
}

func GPTSearchModeName(lang string) string {
	var res string
	if lang == "ja" {
		res = "類似単語検索"
	} else if lang == "en" {
		res = "Similar word search"
	} else {
		res = "類似単語検索"
	}
	return res
}