package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	RefreshToken struct {
		gorm.Model
		Id        uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false;not null;"`
		UserGUID  uuid.UUID `gorm:"type:uuid;not null;"`
		IpAddress string    `gorm:"not null;"`
		Token     string    `gorm:"not null;"`
	}
)

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	r.Id = uuid.New()

	return nil
}
