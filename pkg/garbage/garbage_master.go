package garbage

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"go-academy-presentation/pkg/db"
	"net/http"
	"os"
	"path"

	_ "github.com/go-sql-driver/mysql"
)

const (
	url = "https://www.opendata.metro.tokyo.lg.jp/setagaya/131121_setagayaku_garbage_separate.csv"
)

type GarbageMaster struct {
	Id         int
	PublicCode string
	GarbageId  string
	PublicName string
	District   string
	Item       string
	ItemKana   string
	ItemEng    string
	Classify   string
	Note       string
	Remarks    string
	LargeFee   string
}

func Migrate() {
	db, err := db.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	MigrateMasterGarbage(db)

	// delete download file
	os.Remove(path.Base(url))
}

func MigrateMasterGarbage(db *sql.DB) {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	r := csv.NewReader(res.Body)

	// headerを読み飛ばす
	_, err = r.Read()
	if err != nil {
		panic(err.Error())
	}

	// TODO: 余裕があればトランザクション対応する
	// tx, err := db.Begin()
	// if err != nil {
	// 	panic(err.Error())
	// }

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		var existingGarbageId string
		err = db.QueryRow("SELECT garbage_id FROM garbage_masters WHERE garbage_id = ?", record[1]).Scan(&existingGarbageId)

		switch {
		case err == sql.ErrNoRows:
			gm := &GarbageMaster{
				PublicCode: record[0],
				GarbageId:  record[1],
				PublicName: record[2],
				District:   record[3],
				Item:       record[4],
				ItemKana:   record[5],
				ItemEng:    record[6],
				Classify:   record[7],
				Note:       record[8],
				Remarks:    record[9],
				LargeFee:   record[10],
			}

			// garbage_idが存在しない場合新しい行を追加
			result, err := CreateGm(db, gm)
			if err != nil {
				panic(err.Error())
			}

			garbageMasterId, err := result.LastInsertId()
			if err != nil {
				panic(err.Error())
			}

			MigrateGarbageItem(db, *gm, garbageMasterId)
		case err != nil:
			panic(err.Error())
		default:
			// garbage_idが存在する場合なにもしない
			fmt.Println("garbage_id already exists")
		}
	}
	// tx.Commit()
}

func CreateGm(db *sql.DB, gm *GarbageMaster) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO garbage_masters (public_code, garbage_id, public_name, district, item, item_kana, item_eng, classify, note, remarks, large_fee) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", gm.PublicCode, gm.GarbageId, gm.PublicName, gm.District, gm.Item, gm.ItemKana, gm.ItemEng, gm.Classify, gm.Note, gm.Remarks, gm.LargeFee)
	if err != nil {
		return result, err
	}
	return result, nil
}

// func Read(db *sql.DB, id int) (*GarbageMaster, error) {
// 	gm := &GarbageMaster{}
// 	err := db.QueryRow("SELECT id, public_code, garbage_id, public_name, district, item, item_kana, item_eng, classify, note, remarks, large_fee FROM garbage_masters WHERE id = ?", id).Scan(&gm.Id, &gm.PublicCode, &gm.GarbageId, &gm.PublicName, &gm.District, &gm.Item, &gm.ItemKana, &gm.ItemEng, &gm.Classify, &gm.Note, &gm.Remarks, &gm.LargeFee)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gm, nil
// }

// func Update(db *sql.DB, gm *GarbageMaster) error {
// 	_, err := db.Exec("UPDATE garbage_masters SET public_code = ?, garbage_id = ?, public_name = ?, district = ?, item = ?, item_kana = ?, item_eng = ?, classify = ?, note = ?, remarks = ?, large_fee = ? WHERE id = ?", gm.PublicCode, gm.GarbageId, gm.PublicName, gm.District, gm.Item, gm.ItemKana, gm.ItemEng, gm.Classify, gm.Note, gm.Remarks, gm.LargeFee, gm.Id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func Delete(db *sql.DB, id int) error {
// 	_, err := db.Exec("DELETE FROM garbage_masters WHERE id = ?", id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
