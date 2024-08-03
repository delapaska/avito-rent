package flat

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/delapaska/avito-rent/models"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) CreateFlat(flat models.Flat) (models.Flat, error) {

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v\n", err)
		return models.Flat{}, err
	}
	defer tx.Rollback()

	queryInsert := `
		INSERT INTO flat (house_id, price, rooms, status)
		VALUES ($1, $2, $3, 'created')
		RETURNING id, house_id, price, rooms, status`

	var insertedFlat models.Flat
	err = tx.QueryRow(queryInsert, flat.House_id, flat.Price, flat.Rooms).Scan(
		&insertedFlat.Id,
		&insertedFlat.House_id,
		&insertedFlat.Price,
		&insertedFlat.Rooms,
		&insertedFlat.Status,
	)
	if err != nil {
		log.Printf("Error executing insert query: %v\n", err)
		return models.Flat{}, err
	}

	queryUpdateHouse := `
		UPDATE house
		SET updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err = tx.Exec(queryUpdateHouse, flat.House_id)
	if err != nil {
		log.Printf("Error executing update query: %v\n", err)
		return models.Flat{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return models.Flat{}, err
	}

	return insertedFlat, nil
}

func (s *Store) UpdateFlatStatus(userID uuid.UUID, flat models.UpdateStatusPayload) (models.Flat, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v\n", err)
		return models.Flat{}, err
	}
	defer tx.Rollback()

	var currentStatus models.FlatStatus
	var currentModeratorID uuid.UUID
	queryGetStatus := `
		SELECT status, moderator_id
		FROM flat
		WHERE id = $1
		FOR UPDATE`
	err = tx.QueryRow(queryGetStatus, flat.Id).Scan(&currentStatus, &currentModeratorID)
	if err != nil {
		log.Printf("Error fetching current status: %v\n", err)
		return models.Flat{}, err
	}

	if flat.Status == models.StatusOnModeration {
		if currentStatus != models.StatusCreated {
			return models.Flat{}, fmt.Errorf("cannot put flat into moderation from status %s", currentStatus)
		}
		updateModeratorQuery := `
			UPDATE flat
			SET status = $1, moderator_id = $2
			WHERE id = $3`
		_, err = tx.Exec(updateModeratorQuery, flat.Status, userID, flat.Id)
		if err != nil {
			log.Printf("Error executing update query: %v\n", err)
			return models.Flat{}, err
		}
	} else {
		if currentStatus != models.StatusOnModeration {
			return models.Flat{}, fmt.Errorf("flat must be in status 'on moderation' to be approved or declined")
		}
		if flat.Status != models.StatusApproved && flat.Status != models.StatusDeclined {
			return models.Flat{}, fmt.Errorf("invalid status change: %s", flat.Status)
		}
		if currentModeratorID != userID {
			return models.Flat{}, fmt.Errorf("only the assigned moderator can change the status")
		}
		queryUpdate := `
			UPDATE flat
			SET status = $1
			WHERE id = $2`
		_, err = tx.Exec(queryUpdate, flat.Status, flat.Id)
		if err != nil {
			log.Printf("Error executing update query: %v\n", err)
			return models.Flat{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return models.Flat{}, err
	}

	var updatedFlat models.Flat
	queryGetUpdatedFlat := `
		SELECT id, house_id, price, rooms, status
		FROM flat
		WHERE id = $1`
	err = s.db.QueryRow(queryGetUpdatedFlat, flat.Id).Scan(
		&updatedFlat.Id,
		&updatedFlat.House_id,
		&updatedFlat.Price,
		&updatedFlat.Rooms,
		&updatedFlat.Status,
	)
	if err != nil {
		log.Printf("Error fetching updated flat: %v\n", err)
		return models.Flat{}, err
	}

	return updatedFlat, nil
}
