package models

import "time"

type GormModel struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
}
