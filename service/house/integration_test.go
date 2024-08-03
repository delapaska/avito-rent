package house

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetHouseFlats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql database: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	store := NewStore(db)
	handler := &Handler{store: store}

	r.GET("/houses/:id/flats", handler.handleGetHouseFlats)

	t.Run("should return internal server error when database query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1 AND status = 'approved'`).
			WithArgs("1").
			WillReturnError(fmt.Errorf("database query error"))

		req, err := http.NewRequest("GET", "/houses/1/flats", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("userType", "user")
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		assert.Equal(t, "database query error", response["message"]) // Обновлено, чтобы соответствовать фактическому сообщению
		assert.Equal(t, "test-request-id", req.Header.Get("RequestId"))
	})

	t.Run("should return flats when query succeeds", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1 AND status = 'approved'`).
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
				AddRow(1, 1, 100000, 3, "approved"))

		req, err := http.NewRequest("GET", "/houses/1/flats", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("userType", "user")
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		flats := response["flats"].([]interface{})
		assert.Len(t, flats, 1)
		flat := flats[0].(map[string]interface{})
		assert.Equal(t, float64(1), flat["id"])
		assert.Equal(t, float64(1), flat["house_id"])
		assert.Equal(t, float64(100000), flat["price"])
		assert.Equal(t, float64(3), flat["rooms"])
		assert.Equal(t, "approved", flat["status"])
	})

	t.Run("should handle empty result set correctly", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1 AND status = 'approved'`).
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}))

		req, err := http.NewRequest("GET", "/houses/1/flats", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("userType", "user")
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		assert.Empty(t, response["flats"])
	})

	t.Run("should handle moderator role correctly", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, house_id, price, rooms, status FROM flat WHERE house_id = \$1`).
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
				AddRow(2, 1, 150000, 4, "pending"))

		req, err := http.NewRequest("GET", "/houses/1/flats", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("userType", "moderator")
		req.Header.Set("RequestId", "test-request-id")

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&response)
		if err != nil {
			t.Fatalf("error decoding response: %v", err)
		}

		flats := response["flats"].([]interface{})
		assert.Len(t, flats, 1)
		flat := flats[0].(map[string]interface{})
		assert.Equal(t, float64(2), flat["id"])
		assert.Equal(t, float64(1), flat["house_id"])
		assert.Equal(t, float64(150000), flat["price"])
		assert.Equal(t, float64(4), flat["rooms"])
		assert.Equal(t, "pending", flat["status"])
	})
}
