package controllers

import (
	"log"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories"
	"github.com/A1essandr0/umf-server/internal/utils"
	"github.com/go-redis/redis/v8"
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

func (c *LinksController) CreateLink(payload models.RequestBody, hashLength int, ip string) (*models.ResponseBody, error) {
	var key string
	var response = &models.ResponseBody{Link: "", OriginalUrl: payload.Url}

	if payload.Alias == "" {
		for {
			key = utils.CreateUUID(hashLength)
			_, err := c.KV.GetKVStoreRecord(key)
			if err != redis.Nil && err != nil { 
				log.Printf("Error while checking uuid key existence: %s; %+v", key, err)
				return nil, err
			}
			if err == redis.Nil { break }
		}
		err := c.KV.CreateKVStoreRecord(key, payload.Url)
		if err != nil { 
			log.Printf("Error while creating uuid key: %s; %v", key, err)
			return nil, err
		}
		response.Link = key
		
	} else {
		_, err := c.KV.GetKVStoreRecord(payload.Alias)
		if err != nil && err != redis.Nil {
			log.Printf("Error while checking alias key existence: %s; %+v", key, err)
			return nil, err
		}
		if err != redis.Nil {
			log.Printf("Alias key already exists: %s", payload.Alias)
			return nil, &models.KeyAlreadyExists{}
		}

		err = c.KV.CreateKVStoreRecord(payload.Alias, payload.Url)
		if err != nil { 
			log.Printf("Error while creating alias key: %s; %v", key, err)
			return nil, err
		}
		response.Link = payload.Alias
	}

	c.DB.CreateNewLinkEvent(response.Link, payload.Url, ip)

	return response, nil
}