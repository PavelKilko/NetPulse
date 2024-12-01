package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	Address string `json:"address"`
	GroupID uint   `json:"group_id"`
}
