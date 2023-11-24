package controllers

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
)

type RecordsController struct {
	DB repositories.DBStore	
}

func NewRecordsController(db repositories.DBStore) *RecordsController {
	return &RecordsController{DB: db}
}


func (c *RecordsController) GetRecords(ip string) []models.NewLinkEvent {
	return c.DB.GetNewLinkEvents(ip)
}