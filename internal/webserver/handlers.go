package webserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/utils"
	"github.com/go-redis/redis/v8"
)


func (router *Router) CreateLink(w http.ResponseWriter, r *http.Request) {
	var payload models.RequestBody
	userIP, _ := utils.GetIP(r)
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Bad request; %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("received: %+v\n", payload)

	response, err := router.LinksController.CreateLink(payload, router.Config.HASH_LENGTH, userIP)
	if err != nil && !errors.Is(err, &models.KeyAlreadyExists{}) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	if err != nil && errors.Is(err, &models.KeyAlreadyExists{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	marshaled, err := json.Marshal(response)
	if err == nil {
		w.Write(marshaled)
		log.Printf("responded: %+v\n", *response)
	} else {
		log.Printf("something went wrong while encoding json: %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (router *Router) GetLink(w http.ResponseWriter, r *http.Request) {
	fields := r.Context().Value(struct{}{}).([]string)
	link := fields[0]
	userIP, _ := utils.GetIP(r)

	value, err := router.LinksController.GetLink(link, userIP)
	if err != nil && err != redis.Nil {
		log.Printf("Error while getting value for the key: %s;\n %v\n", link, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == redis.Nil {
		log.Printf("Key %s not found\n", link)
		http.NotFound(w, r)
		return
	}

	log.Printf("Got url: %s, redirecting", value)
	http.Redirect(w, r, value, http.StatusFound)
}


func (router *Router) GetRecords(w http.ResponseWriter, r *http.Request) {
	userIP, _ := utils.GetIP(r)
	log.Printf("Getting %d records\n", router.Config.DEFAULT_RECORDS_AMOUNT_TO_GET)	
	records := router.RecordsController.GetRecords(userIP)

	results := &models.RecordsResponse{Count: len(records), IP: userIP}
	for _, record := range records {
		result := models.RecordResponse{
			Shorturl: record.Key,
			Longurl: record.Value,
			CreatedAt: record.CreatedAt.String(),
		}
		results.Records = append(results.Records, &result)
		log.Printf("got record: %v", result)
	}

	marshaled, err := json.Marshal(results)
	if err == nil {
		w.Write(marshaled)
	}
}