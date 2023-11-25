package db

import "github.com/A1essandr0/umf-server/internal/models"

type MockDB struct {}

func (m *MockDB) CreateClickEvent(link, value, ip string) {}

func (m *MockDB) CreateNewLinkEvent(link, url, ip string) {}

func (m *MockDB) GetNewLinkEvents(ip string) []models.NewLinkEvent {
	return []models.NewLinkEvent{}
}
