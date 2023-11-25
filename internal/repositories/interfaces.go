package repositories

import (
	"github.com/A1essandr0/umf-server/internal/models"
)

type KeyValueStore interface {
	CreateKVStoreRecord(key, value string) error
	GetKVStoreRecord(key string) (string, error)
}

type DBStore interface {
	CreateClickEvent(link, value, ip string)
	CreateNewLinkEvent(link, url, ip string)
	GetNewLinkEvents(ip string) []models.NewLinkEvent
}