package storage

import (
	"go-echo-template/internal/storage/auth"
	"go-echo-template/internal/storage/user"
)

type Storage struct {
	User user.UserRepository
	Auth auth.AuthRepository
}

func NewStorage(user user.UserRepository, auth auth.AuthRepository) *Storage {
	return &Storage{
		User: user,
		Auth: auth,
	}
}
