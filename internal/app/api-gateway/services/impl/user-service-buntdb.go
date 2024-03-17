package impl

import (
	"context"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/services"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/shared"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// UserServiceBuntDB act as a placeholders for demo purpose of how to create the implementation for this service
// following username / email / password / authorities are availble:
//
//	  admin@localhost / admin / admin / ROLE_ADMIN, ROLE_USER
//		 user@localhost / user / user / ROLE_USER
type UserServiceBuntDB struct {
	repo services.UserRepository
}

// NewUserServiceBuntDB create the single-ton instance for this service
func NewUserServiceBuntDB(repo services.UserRepository) services.UserService {
	return &UserServiceBuntDB{
		repo: repo,
	}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceBuntDB) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceBuntDB) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GetUserByUsername return user information based on login username
func (svc *UserServiceBuntDB) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	res, err := svc.repo.FindByLogin(login)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceBuntDB) IsValidToken(t *jwt.Token, login string) bool {
	if !t.Valid {
		slog.Error("token is invalid")
		return false
	}
	sub, err := t.Claims.GetSubject()
	if err != nil {
		slog.Error("failed to get subject from token", "err", err)
		return false
	}
	if sub != login {
		slog.Error("subject is not same as login", "sub", sub, "login", login)
		return false
	}
	return true
}

// RegisterAccount register for a new account
func (svc *UserServiceBuntDB) RegisterAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return svc.repo.Save(account)
}

// SaveAccount save the current account
func (svc *UserServiceBuntDB) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return svc.repo.Save(account)
}
