package usecases

import (
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
)


type RecordsUseCase struct {
	DB repositories.DBStore
	KVStore repositories.KeyValueStore
}

func NewRecordsUseCase(db repositories.DBStore, kv repositories.KeyValueStore) RecordsUseCase {
	return RecordsUseCase{DB: db, KVStore: kv}
}

func (uc *RecordsUseCase) GetRecords() []*models.RecordsResponse {
	return nil
}