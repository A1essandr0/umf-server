package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/utils"
	"github.com/go-redis/redis/v8"
)


func (router *Router) CreateLink(w http.ResponseWriter, r *http.Request) {
	var key string
	var payload models.RequestBody
	var response = &models.ResponseBody{Link: "", OriginalUrl: ""}

	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Bad request; %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// log.Println(r.Header)
	log.Printf("received: %+v\n", payload)
	response.OriginalUrl = payload.Url



	if payload.Alias == "" {
		for {
			key = utils.CreateUUID(router.Config.HASH_LENGTH)
			_, err := router.LinksController.KV.GetKVStoreRecord(key)
			// _, err := router.LinksController.KV.GetKVStoreRecord(key)
			if err != redis.Nil && err != nil { 
				log.Printf("Error while checking uuid key existence: %s; %v", key, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 
			}
			if err == redis.Nil { break }
		}
		err := router.LinksController.KV.CreateKVStoreRecord(key, payload.Url)
		// err := router.LinksController.KV.CreateKVStoreRecord(key, payload.Url)
		if err != nil { 
			log.Printf("Error while creating uuid key: %s; %v", key, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		response.Link = key
		
	} else {
		_, err := router.LinksController.KV.GetKVStoreRecord(payload.Alias)
		// _, err := router.LinksController.KV.GetKVStoreRecord(payload.Alias)
		if err != nil && err != redis.Nil {
			log.Printf("Error while checking alias key existence: %s; %v", key, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		if err != redis.Nil {
			log.Printf("Alias key already exists: %s", payload.Alias)
			http.Error(w, "Alias key already exists", http.StatusConflict)
			return
		}

		err = router.LinksController.KV.CreateKVStoreRecord(payload.Alias, payload.Url)
		// err = router.LinksController.KV.CreateKVStoreRecord(payload.Alias, payload.Url)
		if err != nil { 
			log.Printf("Error while creating alias key: %s; %v", key, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		response.Link = payload.Alias
	}

	userIP, _ := utils.GetIP(r)
	router.RecordsController.DB.CreateNewLinkEvent(response.Link, payload.Url, userIP)
	// router.RecordsController.DB.CreateNewLinkEvent(response.Link, payload.Url, userIP)



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