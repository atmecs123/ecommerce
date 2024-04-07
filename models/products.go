package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func CreateProduct(db *gorm.DB, product *Product) error {
	err := db.Create(product).Error
	if err != nil {
		fmt.Println("Error adding the product to database", err)
	}
	return err
}

func UpdateProduct(db *gorm.DB, product *Product) error {
	err := db.Save(product).Error
	if err != nil {
		fmt.Println("Error adding the product to database", err)
	}
	return err
}
func DeleteProduct(db *gorm.DB, product *Product, id int) error {

	err := db.Where("id = ?", id).Delete(product).Error
	if err != nil {
		fmt.Println("Error deleting the product from database", err)
	}
	return err
}
func GetProduct(db *gorm.DB, product *Product, id int) error {
	err := db.Where("id = ?", id).First(product).Error
	if err != nil {
		fmt.Println("Error getting the product from database", err)
	}
	return err
}

func GetProducts(db *gorm.DB, product *[]Product) error {
	err := db.Find(product).Error
	if err != nil {
		fmt.Println("Error getting the products to database", err)
	}
	return err
}