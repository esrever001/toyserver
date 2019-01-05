package db

import (
	"github.com/esrever001/toyserver/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var schemas = []interface{}{
	&Events{},
}

type Database struct {
	Type     string
	Filename string
	Database *gorm.DB
}

func (db *Database) Init() {
	var err error
	defer utils.HandleError(err)

	db.Database, err = gorm.Open(db.Type, db.Filename)
	if err != nil {
		return
	}

	db.Database.AutoMigrate(schemas...)
}
