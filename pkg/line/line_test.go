package line

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeLanguage(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		want := "en"
		got := ChangeLanguage("ja")
		assert.Equal(t, want, got)
	})
	t.Run("異常系", func(t *testing.T) {
		want := "en"
		got := ChangeLanguage("en")
		assert.NotEqual(t, want, got)
	})
}

func TestChangeSearchMode(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		want := "sql"
		got := ChangeSearchMode("gpt")
		assert.Equal(t, want, got)
	})
	t.Run("異常系", func(t *testing.T) {
		want := "gpt"
		got := ChangeSearchMode("gpt")
		assert.NotEqual(t, want, got)
	})
}
