package flat

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delapaska/avito-rent/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateFlat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql database: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	t.Run("should return error when starting transaction fails", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(fmt.Errorf("transaction start error"))

		flat := models.Flat{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}

		_, err := store.CreateFlat(flat)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		assert.Equal(t, "transaction start error", err.Error())
	})

	t.Run("should return error when insert query fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO flat \(house_id, price, rooms, status\) VALUES \(\$1, \$2, \$3, 'created'\) RETURNING id, house_id, price, rooms, status`).
			WithArgs(1, 100000, 3).
			WillReturnError(fmt.Errorf("insert query error"))
		mock.ExpectRollback()

		flat := models.Flat{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}

		_, err := store.CreateFlat(flat)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		assert.Equal(t, "insert query error", err.Error())
	})

	t.Run("should return error when update query fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO flat \(house_id, price, rooms, status\) VALUES \(\$1, \$2, \$3, 'created'\) RETURNING id, house_id, price, rooms, status`).
			WithArgs(1, 100000, 3).
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
				AddRow(1, 1, 100000, 3, "created"))
		mock.ExpectExec(`UPDATE house SET updated_at = CURRENT_TIMESTAMP WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(fmt.Errorf("update query error"))
		mock.ExpectRollback()

		flat := models.Flat{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}

		_, err := store.CreateFlat(flat)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		assert.Equal(t, "update query error", err.Error())
	})

	t.Run("should successfully create flat and update house", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO flat \(house_id, price, rooms, status\) VALUES \(\$1, \$2, \$3, 'created'\) RETURNING id, house_id, price, rooms, status`).
			WithArgs(1, 100000, 3).
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
				AddRow(1, 1, 100000, 3, "created"))
		mock.ExpectExec(`UPDATE house SET updated_at = CURRENT_TIMESTAMP WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		flat := models.Flat{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}

		createdFlat, err := store.CreateFlat(flat)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedFlat := models.Flat{
			Id:       1,
			House_id: 1,
			Price:    100000,
			Rooms:    3,
			Status:   "created",
		}

		assert.Equal(t, expectedFlat, createdFlat)
	})
}
