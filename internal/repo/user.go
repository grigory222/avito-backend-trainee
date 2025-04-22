package repo

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// AddNewUser user register
func (ur *UserRepository) AddNewUser(email, passwordHash, role string) (models.User, error) {
	sql := "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, email, password_hash, role"
	user := models.User{}
	err := ur.db.QueryRow(sql, email, passwordHash, role).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return models.User{}, models.ErrDBInsert
	}
	return user, nil
}

// GetUserByEmailAndPassword GetUserByEmail user login
func (ur *UserRepository) GetUserByEmailAndPassword(email, passwordHash string) (models.User, error) {
	sql := "SELECT id, email, password_hash, role FROM users WHERE email = $1 AND password_hash = $2"
	user := models.User{}
	err := ur.db.QueryRow(sql, email, passwordHash).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return models.User{}, models.ErrInvalidCredentials
	}
	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (models.User, error) {
	sql := "SELECT id, email, password_hash, role FROM users WHERE email = $1"
	user := models.User{}
	err := ur.db.QueryRow(sql, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return models.User{}, models.ErrDBRead
	}
	return user, nil
}
