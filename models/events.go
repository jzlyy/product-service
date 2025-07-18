package models

import (
	"encoding/json"
	"time"
)

// 事件类型常量
const (
	EventProductCreated  = "product_created"
	EventProductUpdated  = "product_updated"
	EventProductDeleted  = "product_deleted"
	EventCategoryCreated = "category_created"
	EventImageAdded      = "image_added"
	EventAttributeAdded  = "attribute_added"
)

// ProductEvent 商品事件结构
type ProductEvent struct {
	EventID     string           `json:"event_id"`
	EventType   string           `json:"event_type"`
	Timestamp   time.Time        `json:"timestamp"`
	ProductID   int              `json:"product_id"`
	ProductData Product          `json:"product_data,omitempty"`
	CategoryID  int              `json:"category_id,omitempty"`
	ImageData   ProductImage     `json:"image_data,omitempty"`
	Attribute   ProductAttribute `json:"attribute_data,omitempty"`
}

// ToJSON 将事件转换为JSON
func (e *ProductEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
