package model

import "time"

type Category struct {
	ID uint `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"type:varchar(50);not null;unique"`
	//CreatedAt Time `json:"created_at" gorm:"timestamp"`
	//UpdatedAt Time `json:"updated_at" gorm:"timestamp"`
	CreatedAt time.Time `json:"created_at" gorm:"timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"timestamp"`
}