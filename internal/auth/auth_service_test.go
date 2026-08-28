package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"test-product-api/internal/models"
	"test-product-api/internal/repository"
	"test-product-api/internal/test"
	"testing"
)

func setupAuthService(t *testing.T) (*AuthService, func()) {
	db, cleanup := test.SetupTestDBUsers(t)

	userRepo := repository.NewUserRepository(db)
	service := NewAuthService(userRepo, "test-jwt-secret")

	return service, cleanup
}

func TestAuthService_Register(t *testing.T) {
	service, cleanup := setupAuthService(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     models.RegisterRequest
		wantErr bool
	}{
		{
			name: "successful registration",
			req: models.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "123456",
				FullName: "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: models.RegisterRequest{
				Email:    "test@gmail.com",
				Password: "123456",
				FullName: "Another User",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: models.RegisterRequest{
				Email:    "",
				Password: "123456",
				FullName: "Test User",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: models.RegisterRequest{
				Email:    "test2@mail.com",
				Password: "",
				FullName: "Test User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Register(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	service, cleanup := setupAuthService(t)
	defer cleanup()

	registerReq := models.RegisterRequest{
		Email:    "test@gmail.com",
		Password: "123456",
		FullName: "Test User",
	}
	err := service.Register(registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     models.LoginRequest
		wantErr bool
	}{
		{
			name: "successful Login",
			req: models.LoginRequest{
				Email:    "test@gmail.com",
				Password: "123456",
			},
			wantErr: false,
		},
		{
			name: "non-existing email",
			req: models.LoginRequest{
				Email:    "nonexistent@gmail.com",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "wrong password",
			req: models.LoginRequest{
				Email:    "test@gmail.com",
				Password: "worng password",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: models.LoginRequest{
				Email:    "",
				Password: "123456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Login(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthService_Login_ReturnsValidToken(t *testing.T) {
	service, cleanup := setupAuthService(t)
	defer cleanup()

	registerReq := models.RegisterRequest{
		Email:    "test@gmail.com",
		Password: "123456",
		FullName: "Test User",
	}
	err := service.Register(registerReq)
	require.NoError(t, err)

	loginReq := models.LoginRequest{
		Email:    "test@gmail.com",
		Password: "123456",
	}
	token, err := service.Login(loginReq)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ValidateToken(token, "test-jwt-secret")
	require.NoError(t, err)
	assert.Equal(t, float64(1), claims["user_id"])
	assert.Equal(t, "test@gmail.com", claims["email"])
}

func ValidateToken(tokenString string, secret string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
