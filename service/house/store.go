package house

import (
	"database/sql"
	"log"
	"time"

	"github.com/delapaska/avito-rent/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) CreateHouse(house models.House) (models.House, error) {
	log.Println("Create House")
	query := `
		INSERT INTO house (address, year, developer, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, address, year, developer, created_at, updated_at`

	var insertedHouse models.House
	err := s.db.QueryRow(query, house.Address, house.Year, house.Developer).Scan(
		&insertedHouse.Id,
		&insertedHouse.Address,
		&insertedHouse.Year,
		&insertedHouse.Developer,
		&insertedHouse.Created_at,
		&insertedHouse.Updated_at,
	)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return models.House{}, err
	}

	return insertedHouse, nil
}

func (s *Store) GetHouseFlats(houseID string, userRole string) ([]models.Flat, error) {
	var query string
	var args []interface{}

	if userRole == "moderator" {

		query = `
			SELECT id, house_id, price, rooms, status
			FROM flat
			WHERE house_id = $1`
		args = append(args, houseID)
	} else {

		query = `
			SELECT id, house_id, price, rooms, status
			FROM flat
			WHERE house_id = $1 AND status = 'approved'`
		args = append(args, houseID)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var flats []models.Flat
	for rows.Next() {
		var flat models.Flat
		if err := rows.Scan(&flat.Id, &flat.House_id, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		flats = append(flats, flat)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return nil, err
	}

	return flats, nil
}
func (s *Store) AddSubscription(houseID, email string) error {
	_, err := s.db.Exec("INSERT INTO subscriptions (house_id, email, created_at) VALUES ($1, $2, $3)", houseID, email, time.Now())
	return err
}
