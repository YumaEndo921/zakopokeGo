package repository

import "zakopokeGo/internal/domain/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByUserID(userID string) (*model.User, error)
	// FindByMail(mail string) (*model.User, error) // 現状シンプルに保つためコメントアウト、必要なら追加
}

type PokemonRepository interface {
	Create(pokemon *model.Pokemon) error
	FindByUserID(userID uint) ([]*model.Pokemon, error)
	FindByID(id uint) (*model.Pokemon, error)
	Save(pokemon *model.Pokemon) error
	Delete(id uint) error
	CountByUserID(userID uint) (int64, error)
}
