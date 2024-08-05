package flat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateFlat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql database: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	store := NewStore(db)
	handler := &Handler{store: store}

	r.POST("/flats", handler.handleCreateFlat)

	t.Run("should create flat successfully", func(t *testing.T) {
		payload := models.FlatPayload{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}
		marshalled, _ := json.Marshal(payload)
		currentTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO flat \(house_id, price, rooms, status\) VALUES \(\$1, \$2, \$3, 'created'\) RETURNING id, house_id, price, rooms, status`).
			WithArgs(payload.House_id, payload.Price, payload.Rooms).
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).AddRow(1, payload.House_id, payload.Price, payload.Rooms, "created"))
		mock.ExpectExec(`UPDATE house SET updated_at = \$1 WHERE id = \$2`).
			WithArgs(currentTime, payload.House_id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req, err := http.NewRequest("POST", "/flats", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response models.Flat
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		expected := models.Flat{
			Id:       1,
			House_id: payload.House_id,
			Price:    payload.Price,
			Rooms:    payload.Rooms,
			Status:   "created",
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request when JSON is invalid", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/flats", bytes.NewBuffer([]byte("{invalid-json}")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		assert.Equal(t, "invalid character 'i' looking for beginning of object key string", response["message"])
		assert.Equal(t, "test-request-id", req.Header.Get("RequestId"))
	})

	t.Run("should return bad request when payload is invalid", func(t *testing.T) {
		payload := struct {
			House_id string `json:"house_id"`
			Price    int    `json:"price"`
			Rooms    int    `json:"rooms"`
		}{
			House_id: "invalid",
			Price:    100000,
			Rooms:    3,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/flats", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		message, ok := response["message"].(string)
		if !ok {
			t.Fatalf("expected 'message' to be a string, got %T", response["message"])
		}

		assert.Contains(t, message, "json: cannot unmarshal string into Go struct field FlatPayload.house_id of type int")
		assert.Equal(t, "test-request-id", req.Header.Get("RequestId"))
	})

	t.Run("should return bad request with validation errors", func(t *testing.T) {

		payload := struct {
			House_id int `json:"house_id"`
			Price    int `json:"price"`
			Rooms    int `json:"rooms"`
		}{
			Price: 100000,
			Rooms: 3,
		}
		marshalled, _ := json.Marshal(payload)

		validate := validator.New()

		utils.Validate = validate

		req, err := http.NewRequest("POST", "/flats", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		message, ok := response["message"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected 'message' to be a map, got %T", response["message"])
		}

		assert.Contains(t, message["house_id"].(string), "field validation for 'House_id' failed on the 'required' tag")
		assert.Equal(t, "test-request-id", req.Header.Get("RequestId"))
	})

	t.Run("should return internal server error when CreateFlat fails", func(t *testing.T) {
		payload := models.FlatPayload{
			House_id: 1,
			Price:    100000,
			Rooms:    3,
		}
		marshalled, _ := json.Marshal(payload)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO flat \(house_id, price, rooms, status\) VALUES \(\$1, \$2, \$3, 'created'\) RETURNING id, house_id, price, rooms, status`).
			WithArgs(payload.House_id, payload.Price, payload.Rooms).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		req, err := http.NewRequest("POST", "/flats", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		code, ok := response["code"].(float64)
		if !ok {
			t.Fatalf("expected 'code' to be a float64, got %T", response["code"])
		}

		assert.Equal(t, http.StatusInternalServerError, int(code))
		assert.Equal(t, "test-request-id", req.Header.Get("RequestId"))
		assert.Contains(t, response["message"].(string), "database error")
	})
}
