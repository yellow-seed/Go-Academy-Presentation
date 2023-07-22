package garbage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGi(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gi := &GarbageItem{
			GarbageId:       "test",
			GarbageMasterId: 1,
			Category:        "burnable",
		}
		mock.ExpectExec("INSERT INTO garbage_items").WithArgs(gi.GarbageId, gi.GarbageMasterId, gi.Category).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = CreateGi(db, gi)
		if err != nil {
			t.Error("CreateGi is failed")
		}
	},
	)
	t.Run("異常系", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		gi := &GarbageItem{
			GarbageId:       "test",
			GarbageMasterId: 1,
			Category:        "burnable",
		}
		mock.ExpectExec("INSERT INTO garbage_items").WithArgs(gi.GarbageId, gi.GarbageMasterId, gi.Category).WillReturnError(err)

		_, err = CreateGi(db, gi)
		if err == nil {
			t.Error("CreateGi is failed")
		}
	},
	)
}

func TestCategoryName(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		want := "burnable"
		category := CategoryName("可燃ごみ")
		assert.Equal(t, want, category)
	},
	)
	t.Run("異常系", func(t *testing.T) {
		want := ""
		category := CategoryName("test")
		assert.NotEqual(t, want, category)
	},
	)
}
