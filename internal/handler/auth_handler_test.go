package handler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"test-product-api/internal/auth"
	"test-product-api/internal/models"
	"test-product-api/internal/repository"
	"test-product-api/internal/test"
	"testing"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, func()) {
	db, cleanup := test.SetupTestDBUsers(t)
	repo := repository.NewUserRepository(db)
	authService := auth.NewAuthService(repo, "test-jwt-secret")
	handler := NewAuthHandler(authService)
	return handler, cleanup
}

func TestAuthHandler_Register(t *testing.T) {
	handler, cleanup := setupAuthHandler(t)
	defer cleanup()

	tests := []struct {
		name       string
		payload    interface{}
		wantStatus int
	}{
		{
			name: "successful registration",
			payload: models.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "123456",
				FullName: "Test User",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate email",
			payload: models.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "123456",
				FullName: "Another User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invlaid json",
			payload:    `{"email": "test@gmail.com", "password": "123456"`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	handler, cleanup := setupAuthHandler(t)
	defer cleanup()

	registerReq := models.RegisterRequest{
		Email:    "test@gmail.com",
		Password: "123456",
		FullName: "Test User",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	tests := []struct {
		name       string
		payload    interface{}
		wantStatus int
	}{
		{
			name: "successful login",
			payload: models.LoginRequest{
				Email:    "test@gmail.com",
				Password: "123456",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			payload: models.LoginRequest{
				Email:    "test@gmail.com",
				Password: "wrong_password",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existing user",
			payload: models.LoginRequest{
				Email:    "nonexistent@gmail.com",
				Password: "123456",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid json",
			payload:    `{"email": "test@gmail.com", "password": "123456"`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_login_ReturnsToken(t *testing.T) {
	handler, cleanup := setupAuthHandler(t)
	defer cleanup()

	registerReq := models.RegisterRequest{
		Email:    "test@gmail.com",
		Password: "123456",
		FullName: "Test User",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	loginReq := models.LoginRequest{
		Email:    "test@gmail.com",
		Password: "123456",
	}
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	token, ok := response["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)
}
