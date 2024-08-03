package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type FlatStatus string

const (
	StatusCreated      FlatStatus = "created"
	StatusApproved     FlatStatus = "approved"
	StatusDeclined     FlatStatus = "declined"
	StatusOnModeration FlatStatus = "on moderate"
)

var ValidStatuses = map[FlatStatus]bool{
	StatusCreated:      true,
	StatusApproved:     true,
	StatusDeclined:     true,
	StatusOnModeration: true,
}

type DummyStore interface{}

type DummyLoginPayload struct {
	UserType string `json:"userType" validate:"required"`
}
type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

type HouseStore interface {
	CreateHouse(house House) (House, error)
	GetHouseFlats(houseID string, userRole string) ([]Flat, error)
	AddSubscription(houseID, email string) error
}
type House struct {
	Id         int       `json:"id"`
	Address    string    `json:"address"`
	Year       int       `json:"year"`
	Developer  string    `json:"developer"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type HousePayload struct {
	Address   string `json:"address" validate:"required"`
	Year      int    `json:"year" validate:"required"`
	Developer string `json:"developer"`
}
type FlatStore interface {
	CreateFlat(flat Flat) (Flat, error)
	UpdateFlatStatus(userID uuid.UUID, flat UpdateStatusPayload) (Flat, error)
}

type Flat struct {
	Id       int        `json:"id"`
	House_id int        `json:"house_id"`
	Price    int        `json:"price"`
	Rooms    int        `json:"rooms"`
	Status   FlatStatus `json:"status"`
}

type UpdateStatusPayload struct {
	Status FlatStatus `json:"status"  validate:"required"`
	Id     int        `json:"id"  validate:"required"`
}
type FlatPayload struct {
	House_id int `json:"house_id"  validate:"required"`
	Price    int `json:"price"  validate:"required"`
	Rooms    int `json:"rooms"  validate:"required"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id uuid.UUID) (*User, error)
	CreateUser(user User) error
}
type User struct {
	User_id  uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	UserType string    `json:"userType"`
}
type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=16"`
	UserType string `json:"userType"`
}

type LoginUserPayload struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password" validate:"required"`
}

type Subscription struct {
	ID        int       `json:"id"`
	HouseID   string    `json:"house_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type SubscribePayload struct {
	Email string `json:"email" validate:"required"`
}
