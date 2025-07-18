package models

import (
	"time"
)

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	Stock       int       `json:"stock" binding:"required"`
	CategoryID  int       `json:"category_id" binding:"required"`
	SKU         string    `json:"sku"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductDetail struct {
	Product
	CategoryName string             `json:"category_name"`
	Attributes   []ProductAttribute `json:"attributes,omitempty"`
	Images       []ProductImage     `json:"images,omitempty"`
}

type ProductAttribute struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ImageURL  string `json:"image_url" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductFilter struct {
	CategoryID int     `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	Search     string  `form:"search"`
}

type Pagination struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

type ProductResponse struct {
	Products  []ProductDetail `json:"products"`
	Total     int             `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
