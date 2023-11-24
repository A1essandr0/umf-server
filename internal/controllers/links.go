package controllers

import (
	"github.com/A1essandr0/umf-server/internal/repositories"
)


type LinksController struct {
	KV repositories.KeyValueStore
	DB repositories.DBStore
}

func NewLinksController(kv repositories.KeyValueStore, db repositories.DBStore) *LinksController {
	return &LinksController{KV: kv, DB: db}
}

func (c *LinksController) GetLink(link string, ip string) (string, error) {
	value, err := c.KV.GetKVStoreRecord(link)
	if err != nil {
		return "", err
	}

	c.DB.CreateClickEvent(link, value, ip)
	return value, err
}

func (c *LinksController) CreateLink() {

}