package users

import (
	"log"

	"gorm.io/gorm"
)

type repo struct {
	*gorm.DB
}

func NewUserRepo(db *gorm.DB) *repo {
	return &repo{
		DB: db,
	}
}

type IUserRepository interface {
	GetUsers() (out []*UserEntity, err error)
}

func (r *repo) GetUsers() (out []*UserEntity, err error) {
	err = r.DB.Find(&out).Error
	if err != nil {
		log.Println(err)
		return
	}
	return
}
