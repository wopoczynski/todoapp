package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoModel struct {
	gorm.Model
	ID          uuid.UUID `json:"id"          gorm:"type:char(36);primary_key;"`
	Title       string    `json:"title"`
	Order       int       `json:"order"`
	Completed   bool      `json:"completed"`
	Description string    `json:"description"`
}
