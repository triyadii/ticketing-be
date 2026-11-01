package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusAktif    UserStatus = "aktif"
	UserStatusNonaktif UserStatus = "nonaktif"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nama      string     `gorm:"type:varchar(100);not null" json:"nama"`
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Status    UserStatus `gorm:"type:varchar(20);default:'aktif';not null" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Status == "" {
		u.Status = UserStatusAktif
	}
	return
}
