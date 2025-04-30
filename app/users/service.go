package users

import "gorm.io/gorm"

type service struct {
	UserRepo IUserRepository
}

func NewUserService(db *gorm.DB) *service {
	return &service{
		UserRepo: NewUserRepo(db),
	}
}

type IUserService interface {
	GetUsers() (out []*UserEntity, err error)
}

func (s *service) GetUsers() (out []*UserEntity, err error) {
	out, err = s.UserRepo.GetUsers()
	return
}
