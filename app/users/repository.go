package users

import (
	"log"
	"time"

	"github.com/google/uuid"
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
	GetUsers() ([]*UserEntity, error)
	GetUserByID(id uuid.UUID) (*UserEntity, error)
	CreateUser(user *UserEntity) error
	UpdateUser(user *UserEntity) error
	DeleteUser(id uuid.UUID) error
}

func (r *repo) GetUsers() ([]*UserEntity, error) {
	var users []*UserEntity
	err := r.DB.Find(&users).Error
	if err != nil {
		log.Println("GetUsers error:", err)
		return nil, err
	}
	return users, nil
}

func (r *repo) GetUserByID(id uuid.UUID) (*UserEntity, error) {
	var user UserEntity
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		log.Println("GetUserByID error:", err)
		return nil, err
	}
	return &user, nil
}

func (r *repo) CreateUser(user *UserEntity) error {
	err := r.DB.Create(user).Error
	if err != nil {
		log.Println("CreateUser error:", err)
	}
	return err
}

func (r *repo) UpdateUser(user *UserEntity) error {
	user.UpdatedAt = time.Now()
	err := r.DB.Model(user).Updates(user).Error
	if err != nil {
		log.Println("UpdateUser error:", err)
	}
	return err
}

func (r *repo) DeleteUser(id uuid.UUID) error {
	err := r.DB.Delete(&UserEntity{}, "id = ?", id).Error
	if err != nil {
		log.Println("DeleteUser error:", err)
	}
	return err
}
