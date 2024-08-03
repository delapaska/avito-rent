package house

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delapaska/avito-rent/models"

	"github.com/stretchr/testify/assert"
)

func TestGetHouseFlats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql database: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
			AddRow(1, "1", 100000, 3, "approved").
			AddRow(2, "1", 150000, 4, "approved"))

	flats, err := store.GetHouseFlats("1", "moderator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFlats := []models.Flat{
		{Id: 1, House_id: 1, Price: 100000, Rooms: 3, Status: "approved"},
		{Id: 2, House_id: 1, Price: 150000, Rooms: 4, Status: "approved"},
	}
	assert.Equal(t, expectedFlats, flats)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestGetHouseFlatsWithUserRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql database: %v", err)
	}
	defer db.Close()

	store := NewStore(db)

	mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1 AND status = 'approved'`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
			AddRow(1, "1", 100000, 3, "approved"))

	flats, err := store.GetHouseFlats("1", "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFlats := []models.Flat{
		{Id: 1, House_id: 1, Price: 100000, Rooms: 3, Status: "approved"},
	}
	assert.Equal(t, expectedFlats, flats)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}
