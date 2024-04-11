package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func CreateProduct(ctx context.Context, db *gorm.DB, product *Product) error {
	err := db.Create(product).Error
	if err != nil {
		fmt.Println("Error adding the product to database", err)
	}
	return err
}

func UpdateProduct(ctx context.Context, db *gorm.DB, product *Product) error {
	err := db.Save(product).Error
	if err != nil {
		fmt.Println("Error adding the product to database", err)
	}
	return err
}
func DeleteProduct(ctx context.Context, db *gorm.DB, product *Product, id int) error {

	err := db.Where("id = ?", id).Delete(product).Error
	if err != nil {
		fmt.Println("Error deleting the product from database", err)
	}
	return err
}
func GetProduct(ctx context.Context, db *gorm.DB, product *Product, id int) error {
	err := db.Where("id = ?", id).First(product).Error
	if err != nil {
		fmt.Println("Error getting the product from database", err)
	}
	return err
}

func GetProducts(ctx context.Context, db *gorm.DB, product *[]Product, offset, limit int) error {
	fmt.Println("#### Inside db #####")
	var err error
	if offset != 0 || limit != 0 {
		err = db.Limit(limit).Offset(offset).Find(product).Error
		return err
	}
	err = db.Find(product).Error
	if err != nil {
		fmt.Println("Error getting the products to database", err)
	}
	return err
}
