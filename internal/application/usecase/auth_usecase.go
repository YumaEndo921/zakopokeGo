package usecase

import (
	"errors"
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(userID, mail, password string) error
	Login(userID, password string) (*model.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

func (u *authUseCase) Register(userID, mail, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := model.NewUser(userID, mail, string(hashed))
	return u.userRepo.Create(newUser)
}

func (u *authUseCase) Login(userID, password string) (*model.User, error) {
	user, err := u.userRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
