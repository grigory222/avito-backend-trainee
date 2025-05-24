package usecases_test

import (
	"errors"
	"testing"

	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/grigory222/avito-backend-trainee/internal/usecases"
	"github.com/grigory222/avito-backend-trainee/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserService_UserRegister(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userService := usecases.NewUserService(mockRepo)

	hash := func(s string) string {
		return usecases.Hash(s)
	}

	tests := []struct {
		name      string
		email     string
		password  string
		role      string
		mockSetup func()
		wantUser  models.User
		wantErr   error
	}{
		{
			name:     "success_register",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", "test@example.com").
					Return(models.User{}, models.ErrDBInsert).Once()
				mockRepo.On("AddNewUser", "test@example.com", hash("password123"), "user").
					Return(models.User{Email: "test@example.com", Role: "user"}, nil).Once()
			},
			wantUser: models.User{Email: "test@example.com", Role: "user"},
			wantErr:  nil,
		},
		{
			name:     "user_already_exists",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", "test@example.com").
					Return(models.User{Email: "test@example.com"}, nil).Once()
			},
			wantUser: models.User{},
			wantErr:  models.ErrUserAlreadyExists,
		},
		{
			name:     "add_new_user_error",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", "test@example.com").
					Return(models.User{}, models.ErrDBInsert).Once()
				mockRepo.On("AddNewUser", "test@example.com", hash("password123"), "user").
					Return(models.User{}, assert.AnError).Once()
			},
			wantUser: models.User{},
			wantErr:  assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := userService.UserRegister(tt.email, tt.password, tt.role)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.Email, user.Email)
				assert.Equal(t, tt.wantUser.Role, user.Role)
			}

			mockRepo.AssertExpectations(t)
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	}
}

func TestUserService_UserLogin(t *testing.T) {
	// локальная функция хеширования, дублирует логику usecases.hash
	const (
		email    = "test@example.com"
		password = "password123"
	)

	expectedHash := usecases.Hash(password)

	tests := []struct {
		name       string
		setupMocks func(userRepo *mocks.UserRepository)
		wantUser   models.User
		wantErr    error
	}{
		{
			name: "success_login",
			setupMocks: func(userRepo *mocks.UserRepository) {
				userRepo.
					On("GetUserByEmailAndPassword", email, expectedHash).
					Return(models.User{Email: email}, nil).
					Once()
			},
			wantUser: models.User{Email: email},
			wantErr:  nil,
		},
		{
			name: "login_failed",
			setupMocks: func(userRepo *mocks.UserRepository) {
				userRepo.
					On("GetUserByEmailAndPassword", email, expectedHash).
					Return(models.User{}, errors.New("invalid credentials")).
					Once()
			},
			wantUser: models.User{},
			wantErr:  errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			tt.setupMocks(mockRepo)

			svc := usecases.NewUserService(mockRepo)
			user, err := svc.UserLogin(email, password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Equal(t, models.User{}, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
