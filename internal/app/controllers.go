package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/utils"
	"github.com/go-redis/redis/v8"
)


func serve(w http.ResponseWriter, r *http.Request) {
	for _, route := range Routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && r.Method == route.Method {
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.Handler(w, r.WithContext(ctx))
			return 
		}
	}
	http.NotFound(w, r)
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	var payload models.RequestBody
	var response = &models.ResponseBody{Link: "", OriginalUrl: ""}
	var key string
	r.Body = http.MaxBytesReader(w, r.Body, 8192)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Println(r.Header)
	log.Printf("received: %+v\n", payload)
	response.OriginalUrl = payload.Url

	// TODO also return http errors 
	if payload.Alias == "" {
		for {
			key = utils.CreateUUID(Config.HASH_LENGTH)
			_, err := KVClient.GetKVStoreRecord(key)
			if err != redis.Nil && err != nil { 
				log.Printf("Error while checking uuid key existence: %s; %v", key, err)
				return 
			}
			if err == redis.Nil { break }
		}
		err := KVClient.CreateKVStoreRecord(key, payload.Url)
		if err != nil { 
			log.Printf("Error while creating uuid key: %s; %v", key, err)
			return 
		}
		response.Link = key
		

	} else {
		_, err := KVClient.GetKVStoreRecord(payload.Alias)
		if err != nil && err != redis.Nil {
			log.Printf("Error while checking alias key existence: %s; %v", key, err)
			return 
		}
		if err != redis.Nil {
			log.Printf("Alias key already exists: %s", payload.Alias)
			http.Error(w, "Alias key already exists", http.StatusConflict)
			return
		}

		err = KVClient.CreateKVStoreRecord(payload.Alias, payload.Url)
		if err != nil { 
			log.Printf("Error while creating alias key: %s; %v", key, err)
			return 
		}
		response.Link = payload.Alias
	}

	userIP, _ := utils.GetIP(r)
	DB.CreateNewLinkEvent(response.Link, payload.Url, userIP)

	marshaled, err := json.Marshal(response)
	if err == nil {
		w.Write(marshaled)
		log.Printf("responded: %+v\n", *response)
	} else {
		log.Printf("something went wrong while encoding json: %+v\n", err)	
	}
}


func GetLink(w http.ResponseWriter, r *http.Request) {
	fields := r.Context().Value(ctxKey{}).([]string)
	link := fields[0]

	value, err := KVClient.GetKVStoreRecord(link)
	if err != nil && err != redis.Nil {
		log.Printf("Error while getting value for the key: %s;\n %v\n", link, err)
		return
	}
	if err == redis.Nil {
		log.Printf("Key %s not found\n", link)
		http.NotFound(w, r)
		return
	}

	userIP, _ := utils.GetIP(r)
	DB.CreateClickEvent(link, value, userIP)

	log.Printf("Got url: %s, redirecting...", value)
	http.Redirect(w, r, value, http.StatusFound)
}


func GetRecords(w http.ResponseWriter, r *http.Request) {
	userIP, _ := utils.GetIP(r)
	log.Printf("Getting %d records...\n", Config.DEFAULT_RECORDS_AMOUNT_TO_GET)
	
	records := DB.GetNewLinkEvents(userIP)

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