package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// @Description Status of a flat
const (

	// @Description Flat has been created but not yet approved
	StatusCreated string = "created"

	// @Description Flat has been approved and is available
	StatusApproved string = "approved"

	// @Description Flat has been declined and is not available
	StatusDeclined string = "declined"

	// @Description Flat is under moderation and approval is pending
	StatusOnModeration string = "on moderation"
)

type DummyStore interface{}

// @Description Payload for dummy login
// @Param userType body string true "Type of the user" example "client"
type DummyLoginPayload struct {

	// @example "client"
	UserType string `json:"userType" validate:"required,oneof=client moderator"`
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

// @description House представляет собой структуру данных для хранения информации о доме.
// @name House
// @example { "id": 1, "address": "123 Elm Street", "year": 2020, "developer": "XYZ Construction", "created_at": "2024-08-04T00:00:00Z", "updated_at": "2024-08-04T00:00:00Z" }
type House struct {
	// @description Идентификатор дома
	// @example 1
	Id int `json:"id"`

	// @description Адрес дома
	// @example "123 Elm Street"
	Address string `json:"address"`

	// @description Год постройки
	// @example 2020
	Year int `json:"year"`

	// @description Разработчик или строитель дома
	// @example "XYZ Construction"
	Developer string `json:"developer"`

	// @description Дата создания записи
	// @example "2024-08-04T00:00:00Z"
	Created_at time.Time `json:"created_at"`

	// @description Дата последнего обновления записи
	// @example "2024-08-04T00:00:00Z"
	Updated_at time.Time `json:"updated_at"`
}

// @description HousePayload представляет собой структуру данных для создания или обновления информации о доме.
// @name HousePayload
// @example { "address": "123 Elm Street", "year": 2020, "developer": "XYZ Construction" }
type HousePayload struct {
	// @description Адрес дома
	// @example "123 Elm Street"
	Address string `json:"address" validate:"required"`

	// @description Год постройки
	// @example 2020
	Year int `json:"year" validate:"required"`

	// @description Разработчик или строитель дома
	// @example "XYZ Construction"
	Developer string `json:"developer"`
}
type FlatStore interface {
	CreateFlat(flat Flat) (Flat, error)
	UpdateFlatStatus(userID uuid.UUID, flat UpdateStatusPayload) (Flat, error)
}

// @Description Represents a flat in the system

// @Name Flat
// @Example { "id": 1, "house_id": 101, "price": 1200, "rooms": 3, "status": "created" }
type Flat struct {
	// @Description Unique identifier for the flat
	// @Example 1
	Id int `json:"id"`
	// @Description Unique identifier for the house to which the flat belongs
	// @Example 101
	House_id int `json:"house_id"`
	// @Description Price of the flat
	// @Example 1200
	Price int `json:"price"`
	// @Description Number of rooms in the flat
	// @Example 3
	Rooms int `json:"rooms"`
	// @Description Status of the flat
	// @Example "created"
	Status string `json:"status"`
}

// @Description Payload for updating the status of a flat

// @Name UpdateStatusPayload
// @Example { "status": "approved", "id": 1 }
type UpdateStatusPayload struct {
	// @Description Status to update the flat to
	// @Enum created,approved,declined,on moderation
	// @Example "approved"
	Status string `json:"status" validate:"required oneof=created on moderation approved declined"`
	// @Description Unique identifier of the flat to update
	// @Example 1
	Id int `json:"id" validate:"required"`
}

// @Description Payload for creating a new flat

// @Name FlatPayload
// @Example { "house_id": 101, "price": 1200, "rooms": 3 }
type FlatPayload struct {
	// @Description Unique identifier of the house to which the flat belongs
	// @Example 101
	House_id int `json:"house_id" validate:"required"`
	// @Description Price of the flat
	// @Example 1200
	Price int `json:"price" validate:"required"`
	// @Description Number of rooms in the flat
	// @Example 3
	Rooms int `json:"rooms" validate:"required"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id uuid.UUID) (*User, error)
	CreateUser(user User) error
}

// @Description User information

// @Name User
// @Example { "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6", "email": "user@example.com", "password": "password123", "userType": "client" }
type User struct {
	// @Description Unique identifier of the user
	// @Example "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	User_id uuid.UUID `json:"user_id"`
	// @Description Email address of the user
	// @Example "user@example.com"
	Email string `json:"email"`
	// @Description Password of the user
	// @Example "password123"
	Password string `json:"password"`
	// @Description Type of the user (e.g., client, moderator)
	// @Example "client"
	UserType string `json:"userType"`
}

// RegisterUserPayload представляет данные для регистрации нового пользователя
// @Description Payload for user registration
// @Param email body string true "Email address of the user" example "user@example.com"
// @Param password body string true "Password for the user" example "securePassword123"
// @Param userType body string true "Type of the user. Can be 'client' or 'moderator'" example "client"
type RegisterUserPayload struct {
	// @example user@example.com
	Email string `json:"email" validate:"required,email"`
	// @example securePassword123
	Password string `json:"password" validate:"required,min=3,max=16"`
	// @example client
	UserType string `json:"userType" validate:"required,oneof=client moderator"`
}

// LoginUserPayload представляет данные для входа пользователя
// @Description Payload for user login
// @Param id body string true "UUID of the user" example "f47ac10b-58cc-4372-a567-0e02b2c3d479"
// @Param password body string true "Password for the user" example "securePassword123"
type LoginUserPayload struct {
	// @example f47ac10b-58cc-4372-a567-0e02b2c3d479
	ID uuid.UUID `json:"id" validate:"required"`
	// @example securePassword123
	Password string `json:"password" validate:"required"`
}

// @Description Subscription information

// @Name Subscription
// @Example { "id": 1, "house_id": "house_123", "email": "user@example.com", "created_at": "2023-07-21T17:32:28Z" }
type Subscription struct {
	// @Description Unique identifier of the subscription
	// @Example 1
	ID int `json:"id"`

	// @Description Unique identifier of the house
	// @Example "house_123"
	HouseID string `json:"house_id"`

	// @Description Email address of the subscriber
	// @Example "user@example.com"
	Email string `json:"email"`

	// @Description Date and time when the subscription was created
	// @Example "2023-07-21T17:32:28Z"
	CreatedAt time.Time `json:"created_at"`
}

// @Description Payload for subscribing to house updates

// @Name SubscribePayload
// @Example { "email": "user@example.com" }
type SubscribePayload struct {
	// @Description Email address to subscribe to house updates
	// @Example "user@example.com"
	Email string `json:"email" validate:"required"`
}
