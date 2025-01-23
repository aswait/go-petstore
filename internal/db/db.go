package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
)

type DB interface {
	AutoMigrate(dst ...interface{}) error
}

type dbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

func NewDbConf() *dbConf {
	return &dbConf{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
}

func NewDB(dbConf dbConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name,
	)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(5))

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %d timeout", 5)
		case <-ticker.C:
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			sqlDB, err := db.DB()
			if err != nil {
				return nil, err
			}
			err = sqlDB.Ping()
			if err == nil {
				return db, nil
			}
		}
	}
}

func MigrateDB(db DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Category{},
		&models.Pet{},
		&models.PhotoUrl{},
		&models.Order{},
	)
	if err != nil {
		log.Printf("Migration error: %v", err)
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
