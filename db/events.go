package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Events struct {
	gorm.Model
	User  string
	Type  string
	Time  time.Time
	Notes string
	Image string
}
