package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Post struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId uint `json:"user_id" gorm:"not null"`
	CategoryId uint `json:"category_id" gorm:"not null"`
	Category *Category
	Title string `json:"title" gorm:"varchar(50);not null"`
	HeadImg string `json:"head_img"`
	Content string `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"timestamp"`
}

func (post *Post) BeforeCreate(scope *gorm.Scope) error{
	return scope.SetColumn("ID", uuid.NewV4())
}
