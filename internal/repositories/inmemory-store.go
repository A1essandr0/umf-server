package repositories

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"gorm.io/gorm"
)


type MStore struct {
	DB *gorm.DB
	DEFAULT_RECORDS_AMOUNT_TO_GET int
}


func (d *MStore) CreateClickEvent(link, value, ip string) {
	d.DB.Create(&models.ClickEvent{
		Key: link,
		Value: value,
		UserIP: ip,
	})
}

func (d *MStore) CreateNewLinkEvent(link, url, ip string) {
	d.DB.Create(&models.NewLinkEvent{
		Key: link,
		Value: url,
		UserIP: ip,
	})
}

func (d *MStore) GetNewLinkEvents(ip string) []models.NewLinkEvent {
	var records []models.NewLinkEvent
	d.DB.Where(&models.NewLinkEvent{UserIP: ip},
		).Order("created_at desc",
		).Limit(d.DEFAULT_RECORDS_AMOUNT_TO_GET,
		).Find(&records)
	return records
}