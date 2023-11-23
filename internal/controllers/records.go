package controllers

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/usecases"
)

type RecordsController struct {
	LinksUC usecases.LinksUseCase
	RecordsUC usecases.RecordsUseCase
}

func (controller *RecordsController) GetRecords() []*models.RecordsResponse {
	return controller.RecordsUC.GetRecords()
}