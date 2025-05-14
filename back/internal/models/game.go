package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	CreatedBy      string          `json:"createdBy,omitempty"`
	GameCategories []GameCategory  `json:"gameCategories,omitempty" gorm:"foreignKey:GameID"` // has-many relationship
}

type GameCategory struct {
	ID            uint            `gorm:"primaryKey" json:"-"`
	CreatedAt     time.Time       `json:"-"`
	UpdatedAt     time.Time       `json:"-"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"-"`
	GameID        uint            `json:"-"` // Foreing key back to Game
	CategoryTitle string          `json:"categoryTitle,omitempty"`
	Words         WordList        `json:"words,omitempty" gorm:"type:text"`
}
