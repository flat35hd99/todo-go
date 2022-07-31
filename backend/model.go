package main

import "time"

// Costumized gorm.Model
type BaseStruct struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseStruct

	Name string `json:"name"`
	Age  uint   `json:"age"`

	Todos []Todo `json:"todos"`
}

type Todo struct {
	BaseStruct

	Title string `json:"title"`
	Body  string `json:"body"`
	Done  bool   `json:"done"`

	UserID int `json:"user_id"`
}
