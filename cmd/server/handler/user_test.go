package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabriel-ballesteros/voyagr-api/internal/domain"
	user "github.com/gabriel-ballesteros/voyagr-api/internal/user"
	"github.com/gabriel-ballesteros/voyagr-api/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	createReqUser = `{
		"Email": "jd@mail.com",
		"Name": "John Doe",
		"Password": "1"
	}`
	createReqUserConflict = `{
		"Email": "user@mail.com",
		"Name": "John Doe",
		"Password": "1"
	}`
	createReqUserIncomplete = `{
	}`
	changePasswordReq = `{
		"oldPassword": "1234",
		"newPassword": "2"
	}`
	changePasswordReqUnauth = `{
		"oldPassword": "wrong_password",
		"newPassword": "2"
	}`
	updateReqUser = `{
		"Name": "New Name"
	}`

	updateReqUserIncomplete = `{
	}`

	dataUser = domain.User{
		Email:    "user@mail.com",
		Name:     "John Doe",
		Password: "1234",
	}
)

func createServerWithDataUser() *gin.Engine {
	var mockDb map[string]domain.User = map[string]domain.User{"user@mail.com": dataUser}
	service := user.NewMockService(&mockDb)
	userHandler := NewUser(service)
	r := gin.Default()
	userRoutes := r.Group("/api/v1/users")
	{
		userRoutes.GET("/:email", userHandler.Get())
		userRoutes.POST("/create_user", userHandler.Store())
		userRoutes.POST("/:email/reset_password", userHandler.ResetPassword())
		userRoutes.POST("/:email/change_password", userHandler.ChangePassword())
		userRoutes.PATCH("/:email", userHandler.Update())
		userRoutes.DELETE("/:email", userHandler.Delete())
	}

	return r
}

func CreateRequestTestUser(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func TestGetUser_ok(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodGet, "/api/v1/users/user@mail.com", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, dataUser.Name, result.Data.Name)
}

func TestGetUser_not_found(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodGet, "/api/v1/users/nonexisting_user@mail.com", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestCreateUser_ok(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/create_user", createReqUser)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusCreated
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "John Doe", result.Data.Name)
}

func TestCreateUser_conflict(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/create_user", createReqUserConflict)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusConflict
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestCreateUser_bad_request(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/create_user", createReqUserIncomplete)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusBadRequest
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestResetPassword_ok(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/user@mail.com/reset_password", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	data, err := io.ReadAll(rr.Result().Body)
	assert.Nil(t, err)
	assert.Equal(t, "\"Password reseted successfully\"", string(data))
}

func TestResetPassword_not_found(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/nonexistent_user@mail.com/reset_password", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
}

func TestChangePassword_ok(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/user@mail.com/change_password", changePasswordReq)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	data, err := io.ReadAll(rr.Result().Body)
	assert.Nil(t, err)
	assert.Equal(t, "\"Password updated successfully\"", string(data))
}

func TestChangePassword_unauthorized(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPost, "/api/v1/users/user@mail.com/change_password", changePasswordReqUnauth)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusUnauthorized
	assert.Equal(t, expectedCode, rr.Code)
	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
}

func TestUpdateUser_ok(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPatch, "/api/v1/users/user@mail.com", updateReqUser)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "New Name", result.Data.Name)
}

func TestUpdateUser_non_found(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPatch, "/api/v1/users/nonexistent_user@mail.com", updateReqUser)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestUpdateUser_bad_request(t *testing.T) {
	type response struct {
		Data domain.User `json:"data"`
	}
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodPatch, "/api/v1/users/user@mail.com", updateReqUserIncomplete)
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusBadRequest
	assert.Equal(t, expectedCode, rr.Code)
	result := response{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "", result.Data.Name)
}

func TestDeleteUser_ok(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodDelete, "/api/v1/users/user@mail.com", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNoContent
	assert.Equal(t, expectedCode, rr.Code)
}

func TestDeleteUser_not_found(t *testing.T) {
	r := createServerWithDataUser()
	req, rr := CreateRequestTestUser(http.MethodDelete, "/api/v1/users/nonexistent_user@mail.com", "")
	r.ServeHTTP(rr, req)

	expectedCode := http.StatusNotFound
	assert.Equal(t, expectedCode, rr.Code)
	result := web.Error{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.Nil(t, err)
}
