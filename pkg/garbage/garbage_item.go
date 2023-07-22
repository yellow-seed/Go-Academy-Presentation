package garbage

import "database/sql"

type GarbageItem struct {
	Id              int
	GarbageId       string
	GarbageMasterId int
	Category        string
}

func CreateGi(db *sql.DB, gi *GarbageItem) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO garbage_items (garbage_id, garbage_master_id, category) VALUES (?, ?, ?)", gi.GarbageId, gi.GarbageMasterId, gi.Category)
	if err != nil {
		return result, err
	}
	return result, nil
}

func MigrateGarbageItem(db *sql.DB, gm GarbageMaster, garbageMasterId int64) {
	// garbage_itemテーブルにレコードを追加
	gi := &GarbageItem{
		GarbageId:       gm.GarbageId,
		GarbageMasterId: int(garbageMasterId),
		Category:        CategoryName(gm.Classify),
	}
	result, err := CreateGi(db, gi)
	if err != nil {
		panic(err.Error())
	}
	garbageItemId, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	MigrateGarbageItemDetail(db, gm, *gi, garbageItemId)
}

func CategoryName(classify string) string {
	var category string

	if classify == "可燃ごみ" {
		category = "burnable"
	} else if classify == "不燃ごみ" {
		category = "unburnable"
	} else if classify == "資源" {
		category = "recyclable"
	} else if classify == "粗大ごみ" {
		category = "large"
	} else {
		category = "other"
	}
	return category
}
