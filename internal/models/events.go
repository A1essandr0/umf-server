package models

import "gorm.io/gorm"


type NewLinkEvent struct {
	gorm.Model
	Key string
	Value string
	UserIP string
}
type ClickEvent struct {
	gorm.Model
	Key string
	Value string
	UserIP string
}
