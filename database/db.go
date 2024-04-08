package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() (*gorm.DB, error) {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Errorf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	DB, err = ConnectDB()
	return DB, err
}

func retry(maxRetries int, delay time.Duration, fn func() (*gorm.DB, error)) (*gorm.DB, error) {
	fmt.Println("###### Inside retry #####")
	for i := 0; i < maxRetries; i++ {
		gormDb, err := fn()
		if err == nil {
			fmt.Println("###### Retry count ####### ", i)
			return gormDb, err // Success
		}
		fmt.Printf("Attempt %d failed; retrying in %v...\n", i+1, delay)
		time.Sleep(delay)
	}
	return nil, errors.New("maximum retries reached")
}

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp" + "(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?" + "parseTime=true&loc=Local"
	fmt.Println("Full database connection name", dsn)
	gormOpenfn := func() (*gorm.DB, error) {
		gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Error connecting to my sql database", err)
			return gormDB, err
		}
		return gormDB, err
	}
	gormDB, err := retry(4, 2*time.Second, gormOpenfn)
	return gormDB, err
}
