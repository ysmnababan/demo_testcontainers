package users

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Email     string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (UserEntity) TableName() string {
    return "users" // match your actual table name
}