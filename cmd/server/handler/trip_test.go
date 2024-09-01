package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	trip "github.com/gabriel-ballesteros/voyagr-api/internal/trip"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	createReqTrip = `{
		"Name": "Trip Name",
		"Description": "Test",
		"Start": "2024-09-01",
		"End": "2024-11-10",
		"Owner": "user@mail.com",
		"SharedWith": ["user2@mail.com", "user3@mail.com"],
		"Itinerary": []
	}`
	updateReqTrip = `{
		"Name": "Updated Trip Name",
		"Description": "Test",
		"Start": "2024-09-01",
		"End": "2024-11-10",
		"Owner": "user@mail.com",
		"SharedWith": [],
		"Itinerary": []
	}`

	updateReqTripempty = `{
	}`

	createReqTripIncomplete = `{
		"Name": "Trip Name"
	}`

	dataTrip = domain.Trip{
		ID:          "1",
		Name:        "Test trip",
		Description: "Test description",
		Start:       "2024-01-01",
		End:         "2024-02-20",
		Owner:       "user@mail.com",
		SharedWith:  []string{"user2@mail.com", "user3@mail.com"},
		Itinerary:   []domain.ItineraryElement{},
	}
)

func createServerTrip() *gin.Engine {
	var mockDb map[string]domain.Trip = map[string]domain.Trip{}
	service := trip.NewMockService(&mockDb)
	tripHandler := NewTrip(service)
	r := gin.Default()
	tripRoutes := r.Group("/api/v1/trips")
	{
		tripRoutes.GET("", tripHandler.GetAll())
		tripRoutes.GET("/:id", tripHandler.Get())
		tripRoutes.POST("/", tripHandler.Store())
		tripRoutes.PATCH("/:id", tripHandler.Update())
		tripRoutes.DELETE("/:id", tripHandler.Delete())
	}

	return r
}

func createServerWithDataTrip() *gin.Engine {
	var mockDb map[string]domain.Trip = map[string]domain.Trip{"1": dataTrip}
	mockDb["2"] = dataTrip
	service := trip.NewMockService(&mockDb)
	tripHandler := NewTrip(service)
	r := gin.Default()
	tripRoutes := r.Group("/api/v1/trips")
	{
		tripRoutes.GET("", tripHandler.GetAll())
		tripRoutes.GET("/:id", tripHandler.Get())
		tripRoutes.POST("/", tripHandler.Store())
		tripRoutes.PATCH("/:id", tripHandler.Update())
		tripRoutes.DELETE("/:id", tripHandler.Delete())
	}

	return r
}

func CreateRequestTestTrip(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func TestGetAllTrip_ok(t *testing.T) {
	type response struct {
		Data []domain.Trip `json:"data"`
	}
	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodGet, "/api/v1/trips?user_id=user@mail.com", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)

	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, []domain.Trip{dataTrip, dataTrip}, result.Data)
}

func TestGetAllTrip_notFound(t *testing.T) {
	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodGet, "/api/v1/trips?user_id=nonexistentuser@mail.com", "")
	r.ServeHTTP(rr, req)

	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "not_found", result.Code)
}

func TestGetTrip_ok(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}
	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodGet, "/api/v1/trips/1", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, dataTrip.Name, result.Data.Name)
}

func TestGetTrip_notFound(t *testing.T) {
	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodGet, "/api/v1/trips/inexistent_trip_id", "")
	r.ServeHTTP(rr, req)

	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	assert.Nil(t, err)
}

func TestCreateTrip_ok(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	r := createServerTrip()
	req, rr := CreateRequestTestTrip(http.MethodPost, "/api/v1/trips/", createReqTrip)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusCreated
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "Trip Name", result.Data.Name)
}

func TestCreateTrip_bad_request(t *testing.T) {

	r := createServerTrip()
	req, rr := CreateRequestTestTrip(http.MethodPost, "/api/v1/trips/", createReqTripIncomplete)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusBadRequest
	assert.Equal(t, expectedCode, rr.Code)
	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
}

func TestUpdateTrip_ok(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodPatch, "/api/v1/trips/1", updateReqTrip)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "Updated Trip Name", result.Data.Name)
}

func TestUpdateTrip_not_found(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodPatch, "/api/v1/trips/inexistent_trip", updateReqTrip)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestUpdateTrip_bad_request(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodPatch, "/api/v1/trips/1", updateReqTripempty)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusBadRequest
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestDeleteTrip_ok(t *testing.T) {

	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodDelete, "/api/v1/trips/1", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNoContent
	assert.Equal(t, expectedCode, rr.Code)
	assert.Equal(t, 0, len(rr.Body.Bytes()))
}

func TestDeleteTrip_not_found(t *testing.T) {

	r := createServerWithDataTrip()
	req, rr := CreateRequestTestTrip(http.MethodDelete, "/api/v1/trips/inexistent_trip", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
}
