package db

import (
	"Gin_Event_App_Arc/config"
	"Gin_Event_App_Arc/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.DBConfig) *gorm.DB {
	// postgres dns
	fmt.Println("InitDB running")
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SslMode,
	)
	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// migration
	migration(gormDB)
	return gormDB
}

func migration(t *gorm.DB) {
	err := t.AutoMigrate(&models.User{})
	err = t.AutoMigrate(&models.FabricUser{})
	err = t.AutoMigrate(&models.Event{})
	err = t.AutoMigrate(&models.Attendees{})
	err = t.AutoMigrate(&models.Role{})
	err = t.AutoMigrate(&models.UserRole{})
	err = t.AutoMigrate(&models.CategoryGraph{})
	err = t.AutoMigrate(&models.CourseGraph{})
	err = t.AutoMigrate(&models.FileUpload{})
	fmt.Println("Migration func running")
	if err != nil {
		log.Fatal("table is not migrated")
	}
}
