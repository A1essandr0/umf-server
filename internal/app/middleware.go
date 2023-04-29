package app

import (
	"log"
	"net/http"
)


func mockMiddleWare(next http.Handler) http.Handler {
	log.Printf("using mock middleware handler")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {next.ServeHTTP(w, r)})
}


func corsMiddleWare(next http.Handler) http.Handler {
	var allowedCors string
	if Config.DEVELOPMENT_MODE == "development" {
		allowedCors = "*"
	} else if Config.DEVELOPMENT_MODE == "production" {
		allowedCors = Config.PRODUCTION_CORS
	}
	
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("using cors origin: %s", allowedCors)
        w.Header().Set("Access-Control-Allow-Origin", allowedCors)
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

        if r.Method == "OPTIONS" {
            w.WriteHeader(204)
            return
        }
        next.ServeHTTP(w, r)
    })
}