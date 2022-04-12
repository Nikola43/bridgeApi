package controllers

import (
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func Migrate() {
	/*
		// DROP
		GormDB.Migrator().DropTable(&models.User{})
		GormDB.Migrator().DropTable(&models.NodeDB{})

		// CREATE
		GormDB.AutoMigrate(&models.User{})
		GormDB.AutoMigrate(&models.NodeDB{})
	*/
}
