package db

import (
	"github.com/jinzhu/gorm"
)

type Events struct {
	gorm.Model
	User      string
	Type      string
	Timestamp int64
	Notes     string
	Image     string
}
