package usecases

import "github.com/A1essandr0/umf-server/internal/repositories"

type LinksUseCase struct {
	DB repositories.DBStore
	KVStore repositories.KeyValueStore
}

func NewLinksUseCase(db repositories.DBStore, kv repositories.KeyValueStore) LinksUseCase {
	return LinksUseCase{DB: db, KVStore: kv}
}


func (uc *LinksUseCase) CreateLink() {
	uc.CreateLink()
}


func (uc *LinksUseCase) GetLink() {
	uc.GetLink()
}
