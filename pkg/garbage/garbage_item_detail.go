package garbage

import (
	"database/sql"
	"fmt"
	"go-academy-presentation/pkg/deepl"
)

type GarbageItemDetail struct {
	GarbageId             string
	GarbageItemId         int
	LanguageCode          string
	TranslatedName        string
	TranslatedCategory    string
	TranslatedDescription string
}

func CreateGid(db *sql.DB, gid *GarbageItemDetail) error {
	_, err := db.Exec("INSERT INTO garbage_item_details (garbage_id, garbage_item_id, language_code, translated_name, translated_category, translated_description) VALUES (?, ?, ?, ?, ?, ?)", gid.GarbageId, gid.GarbageItemId, gid.LanguageCode, gid.TranslatedName, gid.TranslatedCategory, gid.TranslatedDescription)
	if err != nil {
		return err
	}
	return nil
}

func MigrateGarbageItemDetail(db *sql.DB, gm GarbageMaster, gi GarbageItem, garbageItemId int64) {
	// TODO: 多言語対応しているならそれらの数だけレコードを追加する
	// garbage_item_detailテーブルにレコードを追加
	gid := &GarbageItemDetail{
		GarbageId:             gm.GarbageId,
		GarbageItemId:         int(garbageItemId),
		LanguageCode:          "ja",
		TranslatedName:        gm.Item,
		TranslatedCategory:    TranslatedCategoryName(gi.Category, "ja"),
		TranslatedDescription: TranslatedDescription(gm, gi, "ja"),
	}
	err := CreateGid(db, gid)
	if err != nil {
		panic(err.Error())
	}

	gidEn := &GarbageItemDetail{
		GarbageId:             gm.GarbageId,
		GarbageItemId:         int(garbageItemId),
		LanguageCode:          "en",
		TranslatedName:        gm.ItemEng,
		TranslatedCategory:    TranslatedCategoryName(gi.Category, "en"),
		TranslatedDescription: TranslatedDescription(gm, gi, "en"),
	}
	err = CreateGid(db, gidEn)
	if err != nil {
		panic(err.Error())
	}
}

func TranslatedCategoryName(category string, lang string) string {
	var translatedCategory string

	if lang == "ja" {
		if category == "burnable" {
			translatedCategory = "可燃ごみ"
		} else if category == "unburnable" {
			translatedCategory = "不燃ごみ"
		} else if category == "recyclable" {
			translatedCategory = "資源"
		} else if category == "large" {
			translatedCategory = "粗大ごみ"
		} else {
			translatedCategory = "その他"
		}
	} else if lang == "en" {
		if category == "burnable" {
			translatedCategory = "burnable garbage"
		} else if category == "unburnable" {
			translatedCategory = "non-burnable garbage"
		} else if category == "recyclable" {
			translatedCategory = "recyclable waste"
		} else if category == "large" {
			translatedCategory = "oversized garbage"
		} else {
			translatedCategory = "other"
		}
	} else {
		translatedCategory = "その他"
	}
	return translatedCategory
}

func TranslatedDescription(gm GarbageMaster, gmi GarbageItem, lang string) string {
	var translatedDescription string

	if lang == "ja" {
		if gmi.Category == "other" {
			translatedDescription = gm.Item + "はその他「" + gm.Classify + "」です。\n"
		} else {
			translatedDescription = gm.Item + "は「" + gm.Classify + "」です。\n"
		}

		if gm.Remarks != "" {
			translatedDescription += "[注意事項]\n" + gm.Remarks
		}
	} else if lang == "en" {
		translatedCategory := TranslatedCategoryName(gmi.Category, lang)
		if gmi.Category == "other" {
			other_classify, err := deepl.TranslateText(gm.Classify, "ja", "en")
			if err != nil {
				fmt.Println(err)
			}
			translatedDescription = gm.ItemEng + `is` + translatedCategory + `"` + other_classify + `".\n`
		} else {
			translatedDescription = gm.ItemEng + `is "` + translatedCategory + `".\n`
		}

		if gm.Remarks != "" {
			translatedRemarks, err := deepl.TranslateText(gm.Remarks, "ja", "en")
			if err != nil {
				panic(err.Error())
			}
			translatedDescription += "[Note]\n" + translatedRemarks
		}
	} else {
		translatedDescription = "Other garbage"
	}
	return translatedDescription
}
