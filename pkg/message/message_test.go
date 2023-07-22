package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitchLanguageMessage(t *testing.T) {
	want := "言語を日本語に変更しました。"
	got := SwitchLanguageMessage("ja")
	assert.Equal(t, want, got, "The two words should be the same.")
}

func TestHowToUseMessage(t *testing.T) {
	want := "世田谷区のごみの分類方法を提供するbotです。分別方法を知りたいごみの名前を送信してください。文字一致検索モードとChatGPTを用いた類似単語検索モードがあります。検索モードの切り替えはメニューのボタンからおこなえます。"
	got := HowToUseMessage("ja")
	assert.Equal(t, want, got, "The two words should be the same.")
}

func TestChangeSearchModeMessage(t *testing.T) {
	want := "検索モードを文字一致検索に切り替えました。"
	got := ChangeSearchModeMessage("ja", "sql")
	assert.Equal(t, want, got, "The two words should be the same.")
}

func TestSearchModeName(t *testing.T) {
	want := "文字一致検索"
	got := SearchModeName("sql", "ja")
	assert.Equal(t, want, got, "The two words should be the same.")
}

func TestSQLSearchModeName(t *testing.T) {
	want := "文字一致検索"
	got := SQLSearchModeName("ja")
	assert.Equal(t, want, got, "The two words should be the same.")
}

func TestGPTSearchModeName(t *testing.T) {
	want := "類似単語検索"
	got := GPTSearchModeName("ja")
	assert.Equal(t, want, got, "The two words should be the same.")
}
