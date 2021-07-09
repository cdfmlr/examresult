package model

import (
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Student{},
		&ExamResult{},
	)
	return err
}
