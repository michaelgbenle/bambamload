package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" gorm:"index"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {

	v7, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if m.ID == "" {
		m.ID = v7.String()
	}

	return
}
