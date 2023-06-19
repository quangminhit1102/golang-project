package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"Id"` // Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validator:"required"`
	Price       float32   `json:"price" validator:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      uuid.UUID
}
