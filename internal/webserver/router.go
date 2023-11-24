package webserver

import (
	"context"
	"net/http"
	"regexp"

	"github.com/A1essandr0/umf-server/internal/controllers"
	"github.com/A1essandr0/umf-server/internal/models"
)

type RoutePattern struct {
    Method  string
    Regex   *regexp.Regexp
    Handler http.HandlerFunc
}
func NewRoutePattern(method string, pattern string, handler http.HandlerFunc) RoutePattern {
	return RoutePattern{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type Router struct {
	Config *models.Config
	LinksController *controllers.LinksController
	RecordsController *controllers.RecordsController
}
func NewRouter(
	config *models.Config, 
	linksController *controllers.LinksController, 
	recordsController *controllers.RecordsController,
) *Router {
	return &Router{
		Config: config,
		LinksController: linksController,
		RecordsController: recordsController,
	}
}


func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Routes := []RoutePattern{
		// order is important
		NewRoutePattern("POST", "/create", router.CreateLink),
		NewRoutePattern("GET",  "/records", router.GetRecords),
		NewRoutePattern("GET",  "/([a-zA-Z0-9_-]{2,32})", router.GetLink),
	}
	for _, route := range Routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && r.Method == route.Method {
			ctx := context.WithValue(r.Context(), struct{}{}, matches[1:])
			route.Handler(w, r.WithContext(ctx))
			return 
		}
	}
	http.NotFound(w, r)
}

