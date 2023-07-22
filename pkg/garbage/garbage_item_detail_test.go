package garbage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGid(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gid := &GarbageItemDetail{
			GarbageId:             "1",
			GarbageItemId:         1,
			LanguageCode:          "ja",
			TranslatedName:        "test",
			TranslatedCategory:    "test",
			TranslatedDescription: "test",
		}

		mock.ExpectExec("INSERT INTO garbage_item_details").WithArgs(gid.GarbageId, gid.GarbageItemId, gid.LanguageCode, gid.TranslatedName, gid.TranslatedCategory, gid.TranslatedDescription).WillReturnResult(sqlmock.NewResult(1, 1))

		err = CreateGid(db, gid)
		if err != nil {
			t.Error(err.Error())
		}
	},
	)

	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gid := &GarbageItemDetail{
			GarbageId:             "1",
			GarbageItemId:         1,
			LanguageCode:          "ja",
			TranslatedName:        "test",
			TranslatedCategory:    "test",
			TranslatedDescription: "test",
		}

		mock.ExpectExec("INSERT INTO garbage_item_details").WithArgs(gid.GarbageId, gid.GarbageItemId, gid.LanguageCode, gid.TranslatedName, gid.TranslatedCategory, gid.TranslatedDescription).WillReturnError(err)

		err = CreateGid(db, gid)
		if err == nil {
			t.Error("error should be occured")
		}
	},
	)
}

func TestTranslatedCategoryName(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		want := "可燃ごみ"
		translatedCategory := TranslatedCategoryName("burnable", "ja")
		assert.Equal(t, want, translatedCategory)
	},
	)

	t.Run("異常系", func(t *testing.T) {
		want := ""
		translatedCategory := TranslatedCategoryName("burnable", "en")
		assert.NotEqual(t, want, translatedCategory)
	},
	)
}

func TestTranslatedDescription(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		gm := GarbageMaster{
			GarbageId: "1",
			Item:      "test",
			ItemEng:   "test",
			Classify:  "可燃ごみ",
			Remarks:   "",
		}
		gi := GarbageItem{
			GarbageId: "1",
			Category:  "burnable",
		}
		want := "testは「可燃ごみ」です。\n"
		translatedDescription := TranslatedDescription(gm, gi, "ja")
		assert.Equal(t, want, translatedDescription)
	},
	)

	t.Run("異常系", func(t *testing.T) {
		gm := GarbageMaster{
			GarbageId: "1",
			Item:      "test",
			ItemEng:   "test",
			Classify:  "可燃ごみ",
			Remarks:   "",
		}
		gi := GarbageItem{
			GarbageId: "1",
			Category:  "burnable",
		}

		want := ""
		trancelatedDescription := TranslatedDescription(gm, gi, "en")
		assert.NotEqual(t, want, trancelatedDescription)
	},
	)
}
