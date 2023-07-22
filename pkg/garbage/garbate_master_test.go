package garbage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateGm(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gm := &GarbageMaster{
			PublicCode: "1",
			GarbageId:  "1",
			PublicName: "test",
			District:   "test",
			Item:       "test",
			ItemKana:   "test",
			ItemEng:    "test",
			Classify:   "test",
			Note:       "test",
			Remarks:    "test",
			LargeFee:   "test",
		}
		mock.ExpectExec("INSERT INTO garbage_masters").WithArgs(gm.PublicCode, gm.GarbageId, gm.PublicName, gm.District, gm.Item, gm.ItemKana, gm.ItemEng, gm.Classify, gm.Note, gm.Remarks, gm.LargeFee).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = CreateGm(db, gm)
		if err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}
	})
	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gm := &GarbageMaster{
			PublicCode: "1",
			GarbageId:  "1",
			PublicName: "test",
			District:   "test",
			Item:       "test",
			ItemKana:   "test",
			ItemEng:    "test",
			Classify:   "test",
			Note:       "test",
			Remarks:    "test",
			LargeFee:   "test",
		}

		mock.ExpectExec("INSERT INTO garbage_masters").WithArgs(gm.GarbageId, gm.Item, gm.ItemKana, gm.ItemEng, gm.Classify, gm.Note, gm.Remarks, gm.LargeFee).WillReturnError(err)
		_, err = CreateGm(db, gm)
		if err == nil {
			t.Errorf("error was expected while updating stats: %s", err)
		}
	})
}
