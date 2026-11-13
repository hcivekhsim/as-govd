package tasks

import (
	"context"
)

// Не существует ли уже этот пользователь
// ...

type AuthService struct {
	users *UserRepo
}

func NewAuthService(users *UserRepo) *AuthService {
	return &AuthService{users: users}
}

func (as *AuthService) RegisterService(ctx context.Context, data Register) (int64, error) {
	hash, err := hashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	id, err := as.users.createUser(ctx, data, hash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (as *AuthService) LoginService(ctx context.Context, data Login) (int64, string, string, error) {
	userFound, err := as.users.getUserByLogin(ctx, data.Login)
	if err != nil {
		return 0, "", "", err
	}

	if !checkPasswordHash(data.Password, userFound.Password) {
		return 0, "", "", ErrPassowrdMismatch
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		return 0, "", "", err
	}

	csrfToken, err := generateToken(32)
	if err != nil {
		return 0, "", "", err
	}

	return userFound.UserID, sessionToken, csrfToken, nil
}

// Изменить ошибку "Ничего не найдено"
