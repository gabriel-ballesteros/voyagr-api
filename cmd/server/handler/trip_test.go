package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	trip "github.com/gabriel-ballesteros/voyagr-api/internal/trip"
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
		"SharedWith": [],
		"Itinerary": []
	}`
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

func CreateRequestTestTrip(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func TestCreateTrip_ok(t *testing.T) {
	type response struct {
		Data domain.Trip `json:"data"`
	}

	r := createServerTrip()
	req, rr := CreateRequestTestTrip(http.MethodPost, "/api/v1/trips/", createReqTrip)
	r.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)

	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "Trip Name", result.Data.Name)
}
