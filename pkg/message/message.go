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
		res = `
このbotは世田谷区のごみの分別に回答するbotです。
分別方法を知りたいゴミの名前を入力することで分別方法の說明が返信されます。

🔎マークをおすことで検索方法を切り替えることができます。
検索方法は一致する文字を検索するモードとChatGPTが関連キーワードを検索するモードの2種類です。

現在、日本語と英語での質問に対応しています。
		`
	} else if lang == "en" {
		res = `
This bot is a bot that responds to garbage separation in Setagaya Ward.
By entering the name of the garbage you want to know the separation method, the explanation of the separation method will be sent back.

You can switch the search method by pressing the 🔎 mark.
There are two types of search methods: a mode that searches for matching characters and a mode that ChatGPT searches for related keywords.

Currently, we are responding to questions in Japanese and English.
		`
	}
	return res
}

func InitialStateMessage() string {
	return `
現在
・言語: 日本語
・検索モード: 文字一致検索
です。
	`
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
		res = "ChatGPTによる類似単語検索"
	} else if lang == "en" {
		res = "Similar word search powered by ChatGPT"
	} else {
		res = "ChatGPTによる類似単語検索"
	}
	return res
}

func GPTSearchCautionMessage(lang string) string {
	if lang == "ja" {
		return `
ChatGPTによる検索はベータ機能です。
検索には少し時間がかかり、精度にも注意が必要です。
		`
	} else {
		return `
Search by ChatGPT is a beta feature.
The search takes some time and requires precision.
		`
	}
}

func ForeignerSupportMessage(lang string) string {
	if lang == "ja" {
		return `
正確な分類を知りたい場合は次のリンク先のページから確認してください。
https://www.city.setagaya.lg.jp/mokuji/kurashi/004/001/d00190086.html
		`
	} else {
		return `
If you want to know the exact classification, please check from the following linked page.
https://www.city.setagaya.lg.jp/mokuji/kurashi/004/013/index.html
		`
	}
}

func ErrorMessage(lang string) string {
	if lang == "ja" {
		return `
エラーが発生しました。改めてメッセージを送信してください。
		`
	}  else if lang == "en" {
		return `
An error has occurred. Please send a message again.
		`
	} else {
		return `
エラーが発生しました。改めてメッセージを送信してください。
		`
	}
}