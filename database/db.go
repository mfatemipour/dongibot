package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

func NewDB(filePath string) (*DBHandler, error) {
	var err error
	dbHandler := new(DBHandler)

	dbHandler.DB, err = gorm.Open(sqlite.Open(filePath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	if err := dbHandler.DB.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	if err := dbHandler.DB.AutoMigrate(&Transaction{}); err != nil {
		return nil, err
	}
	if err := dbHandler.DB.AutoMigrate(&Dong{}); err != nil {
		return nil, err
	}
	if err := dbHandler.DB.AutoMigrate(&DongUser{}); err != nil {
		return nil, err
	}
	return dbHandler, nil
}
