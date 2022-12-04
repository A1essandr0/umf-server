package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/config"
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/redisclient"
	"github.com/A1essandr0/umf-server/internal/router"
	"github.com/A1essandr0/umf-server/internal/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ctxKey struct{}

var (
	routes = []router.RoutePattern{
		router.NewRoute("POST", "/create", createLink),
		router.NewRoute("GET",  "/records", getRecords),
		router.NewRoute("GET",  "/([a-zA-Z0-9_-]{2,32})", getLink),
	}
	redisClient *redisclient.RedisClient
	DB *gorm.DB
)

func Run() {
	var redisError, dbError, httpStartError error
	redisClient, redisError = redisclient.NewRedisClient()
	if redisError != nil {
		log.Fatalf("Failed to connect to redis: %s", redisError.Error())
	}
	log.Printf("Got redis instance on %s", config.REDIS_ADDR)

	DB, dbError = gorm.Open(postgres.New(postgres.Config{
		DSN: config.DB_DSN,
		PreferSimpleProtocol: false,
	}), &gorm.Config{})
	if dbError != nil {
		log.Fatalf("Failed to connect to DB: %s", dbError.Error())
	}
	DB.AutoMigrate(&models.NewLinkEvent{}, &models.ClickEvent{})

	server := &http.Server{
		Addr: config.WEB_PORT,
		Handler: corsMiddleWare(http.HandlerFunc(serve)),
	}
	if config.USE_TLS == "true" {
		log.Printf("Starting https server on %s, mode: %s", config.WEB_PORT, config.DEVELOPMENT_MODE)
		httpStartError = server.ListenAndServeTLS(config.CERT_FILE, config.CERT_KEY_FILE)
	} else {
		log.Printf("Starting http server on %s", config.WEB_PORT)
		httpStartError = server.ListenAndServe()
	}

	if errors.Is(httpStartError, http.ErrServerClosed) {
		log.Println("... server stopped")
	} else if httpStartError != nil {
		log.Printf("Error starting server: %s", httpStartError)
	}
}


func serve(w http.ResponseWriter, r *http.Request) {
	for _, route := range routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && r.Method == route.Method {
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.Handler(w, r.WithContext(ctx))
			return 
		}
	}
	http.NotFound(w, r)
}

func createLink(w http.ResponseWriter, r *http.Request) {
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
			key = utils.CreateUUID(config.HASH_LENGTH)
			_, err := redisClient.GetKVStoreRecord(key)
			if err != redis.Nil && err != nil { 
				log.Printf("Error while checking uuid key existence: %s; %v", key, err)
				return 
			}
			if err == redis.Nil { break }
		}
		err := redisClient.CreateKVStoreRecord(key, payload.Url)
		if err != nil { 
			log.Printf("Error while creating uuid key: %s; %v", key, err)
			return 
		}
		response.Link = key
		

	} else {
		_, err := redisClient.GetKVStoreRecord(payload.Alias)
		if err != nil && err != redis.Nil {
			log.Printf("Error while checking alias key existence: %s; %v", key, err)
			return 
		}
		if err != redis.Nil {
			log.Printf("Alias key already exists: %s", payload.Alias)
			http.Error(w, "Alias key already exists", http.StatusConflict)
			return
		}

		err = redisClient.CreateKVStoreRecord(payload.Alias, payload.Url)
		if err != nil { 
			log.Printf("Error while creating alias key: %s; %v", key, err)
			return 
		}
		response.Link = payload.Alias
	}

	userIP, _ := utils.GetIP(r)

	DB.Create(&models.NewLinkEvent{
		Key: response.Link,
		Value: payload.Url,
		UserIP: userIP,
	})

	marshaled, err := json.Marshal(response)
	if err == nil {
		w.Write(marshaled)
		log.Printf("responded: %+v\n", *response)
	} else {
		log.Printf("something went wrong while encoding json: %+v\n", err)	
	}
}


func getLink(w http.ResponseWriter, r *http.Request) {
	fields := r.Context().Value(ctxKey{}).([]string)
	link := fields[0]

	value, err := redisClient.GetKVStoreRecord(link)
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

	DB.Create(&models.ClickEvent{
		Key: link,
		Value: value,
		UserIP: userIP,
	})

	log.Printf("Got url: %s, redirecting...", value)
	http.Redirect(w, r, value, http.StatusFound)
}


func getRecords(w http.ResponseWriter, r *http.Request) {
	userIP, _ := utils.GetIP(r)
	log.Printf("Getting %d records...\n", config.DEFAULT_RECORDS_AMOUNT_TO_GET)
	
	var records []models.NewLinkEvent
	DB.Where(&models.NewLinkEvent{UserIP: userIP},
		).Order("created_at desc",
		).Limit(config.DEFAULT_RECORDS_AMOUNT_TO_GET,
		).Find(&records) 

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