package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"test-product-api/internal/models"
	"test-product-api/internal/test"
	"testing"
)

func setupUserRepo(t *testing.T) (*UserRepository, func()) {
	db, cleanup := test.SetupTestDBUsers(t)
	repo := NewUserRepository(db)

	cleanupRepo := func() {
		cleanup()
	}
	return repo, cleanupRepo
}

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	user := models.User{
		Email:        "test@gmail.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
	}

	err := repo.Create(user)
	require.NoError(t, err)

	saved, err := repo.GetByEmail(user.Email)
	require.NoError(t, err)
	assert.Equal(t, "test@gmail.com", saved.Email)
	assert.Equal(t, "Test User", saved.FullName)
	assert.NotEmpty(t, saved.PasswordHash)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	user := models.User{
		Email:        "test@gmail.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
	}

	err := repo.Create(user)
	require.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing email",
			email:   "test@gmail.com",
			wantErr: false,
		},
		{
			name:    "non-existing email",
			email:   "nonexistent@gmail.com",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	user := models.User{
		Email:        "test@gmail.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
	}

	err := repo.Create(user)
	require.NoError(t, err)

	saved, err := repo.GetByEmail("test@gmail.com")
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing ID",
			id:      saved.ID,
			wantErr: false,
		},
		{
			name:    "non-existing ID",
			id:      999,
			wantErr: true,
		},
		{
			name:    "zero ID",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "test@gmail.com", user.Email)
			}
		})
	}
}

func TestUserRepository_CreateDuplicateEmail(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	user1 := models.User{
		Email:        "test@gmail.com",
		PasswordHash: "hash1",
		FullName:     "User One",
	}

	err := repo.Create(user1)
	require.NoError(t, err)

	user2 := models.User{
		Email:        "test@gmail.com",
		PasswordHash: "hash2",
		FullName:     "User Two",
	}

	err = repo.Create(user2)
	assert.Error(t, err)
}
