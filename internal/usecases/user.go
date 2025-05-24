package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/grigory222/avito-backend-trainee/internal/models"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func hash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (us *UserService) UserRegister(email, password, role string) (models.User, error) {
	_, err := us.userRepo.GetUserByEmail(email)
	if err == nil {
		return models.User{}, models.ErrUserAlreadyExists
	}
	user, err := us.userRepo.AddNewUser(email, hash(password), role)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (us *UserService) UserLogin(email, password string) (models.User, error) {
	return us.userRepo.GetUserByEmailAndPassword(email, hash(password))
}
