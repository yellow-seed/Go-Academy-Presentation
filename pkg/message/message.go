package message

func SwitchLanguageMessage(currentLang string) string {
	var res string
	if currentLang == "ja" {
		res = "言語を日本語に変更しました。"
	} else if currentLang == "en" {
		res = "Changed the language to English."
	}
	return res
}