package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"gymflow/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin(t *testing.T) {
	router := setupTestRouter()

	// 1. Register
	registerReq := user.RegisterRequest{
		Name:           "Test User",
		Email:          "test@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", registerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, registerResp["token"])

	// 2. Login with same credentials
	loginReq := user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/login", loginReq, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp["token"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	router := setupTestRouter()

	registerReq := user.RegisterRequest{
		Name:           "Test User",
		Email:          "duplicate@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	// First registration
	w := makeRequest(t, router, "POST", "/api/v1/auth/register", registerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Second registration with same email
	w = makeRequest(t, router, "POST", "/api/v1/auth/register", registerReq, "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	router := setupTestRouter()

	loginReq := user.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/login", loginReq, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	router := setupTestRouter()

	// 1. Register user
	registerReq := user.RegisterRequest{
		Name:           "Test User",
		Email:          "testuser@example.com",
		Password:       "correctpassword",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", registerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Try to login with wrong password
	loginReq := user.LoginRequest{
		Email:    "testuser@example.com",
		Password: "wrongpassword",
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/login", loginReq, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}