package webserver

import (
	"errors"
	"log"
	"net/http"

	"github.com/A1essandr0/umf-server/internal/controllers"
	"github.com/A1essandr0/umf-server/internal/models"
)

type StdWebServer struct {
	Config *models.Config
	LinksController *controllers.LinksController
	RecordsController *controllers.RecordsController
}

func (s *StdWebServer) Run() {
	router := NewRouter(s.Config, s.LinksController, s.RecordsController)
	server := &http.Server{
		Addr: s.Config.WEB_PORT,
		Handler: s.mockMiddleWare(
			s.corsMiddleWare(router),
		),
	}

	var httpStartError error
	if s.Config.USE_TLS {
		log.Printf("Starting https server on %s, mode: %s", s.Config.WEB_PORT, s.Config.DEVELOPMENT_MODE)
		httpStartError = server.ListenAndServeTLS(s.Config.CERT_FILE, s.Config.CERT_KEY_FILE)
	} else {
		log.Printf("Starting http server on %s, mode: %s", s.Config.WEB_PORT, s.Config.DEVELOPMENT_MODE)
		httpStartError = server.ListenAndServe()
	}

	if errors.Is(httpStartError, http.ErrServerClosed) {
		log.Println("... server stopped")
	} else if httpStartError != nil {
		log.Printf("Error starting server: %s", httpStartError)
	}
}


func (s *StdWebServer) corsMiddleWare(next http.Handler) http.Handler {
	var allowedCors string
	if s.Config.DEVELOPMENT_MODE == "development" {
		allowedCors = "*"
	} else if s.Config.DEVELOPMENT_MODE == "production" {
		allowedCors = s.Config.PRODUCTION_CORS
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

func (s *StdWebServer) mockMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- mock middleware handler ---")
		next.ServeHTTP(w, r)
	})
}