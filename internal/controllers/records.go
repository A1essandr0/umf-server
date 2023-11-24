package controllers

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
)

type RecordsController struct {
	KV repositories.KeyValueStore
	DB repositories.DBStore	
}

func NewRecordsController(kv repositories.KeyValueStore, db repositories.DBStore) *RecordsController {
	return &RecordsController{KV: kv, DB: db}
}


func (c *RecordsController) GetRecords(ip string) []models.NewLinkEvent {
	return c.DB.GetNewLinkEvents(ip)
}