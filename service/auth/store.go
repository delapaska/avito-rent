package auth

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

func scanRowIntoUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)

	err := rows.Scan(
		&user.User_id,
		&user.Email,
		&user.Password,
		&user.UserType,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT user_id, email, password, user_type FROM users WHERE email = $1`

	row := s.db.QueryRow(query, email)

	u := new(models.User)
	err := row.Scan(&u.User_id, &u.Email, &u.Password, &u.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return u, nil
}

func (s *Store) GetUserById(id uuid.UUID) (*models.User, error) {

	query := `SELECT user_id, email, password, user_type FROM users WHERE user_id = $1`
	log.Printf("Executing query: %s with id: %s", query, id)
	row := s.db.QueryRow(query, id)
	u := new(models.User)
	err := row.Scan(&u.User_id, &u.Email, &u.Password, &u.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return u, nil
}

func (s *Store) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (user_id, email, password, user_type)
		VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query, user.User_id, user.Email, user.Password, user.UserType)
	if err != nil {
		return err
	}

	return nil
}
