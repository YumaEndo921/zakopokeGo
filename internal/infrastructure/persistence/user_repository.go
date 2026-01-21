package persistence

import (
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *model.User) error {
	dbUser := FromDomainUser(user)
	result := DB.Create(dbUser)
	if result.Error != nil {
		return result.Error
	}
	// Reflect ID back to domain model if needed, but for now ID is auto-generated
	user.ID = dbUser.ID
	return nil
}

func (r *userRepository) FindByUserID(userID string) (*model.User, error) {
	var dbUser User
	result := DB.Where("user_id = ?", userID).First(&dbUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return dbUser.ToDomain(), nil
}
