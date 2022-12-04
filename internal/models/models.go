package models

import "gorm.io/gorm"


type RequestBody struct {
	Url string
	Alias string
}
type ResponseBody struct {
	Link string
	OriginalUrl string
}

// TODO user profile data
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


type RecordResponse struct {
	Shorturl string
	Longurl string
	CreatedAt string
}
type RecordsResponse struct {
	IP string
	Records []*RecordResponse
	Count int
}