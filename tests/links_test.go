package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/A1essandr0/umf-server/internal/controllers"
	"github.com/A1essandr0/umf-server/internal/models"
	"github.com/A1essandr0/umf-server/internal/repositories/db"
	"github.com/A1essandr0/umf-server/internal/repositories/kv"
	"github.com/A1essandr0/umf-server/internal/webserver"
)

var testConfig = &models.Config{
	DEVELOPMENT_MODE: "development",
	WEB_PORT: "10009",
	PRODUCTION_CORS: "*",
	KVSTORE_TYPE: "mock",
	DBSTORE_TYPE: "sqlite-inmemory",	
	APPLY_MIGRATIONS: true,
	DB_DEBUG_LOG: false,
	DEFAULT_RECORDS_AMOUNT_TO_GET: 5,
	HASH_LENGTH: 8,
}

var alias = "cats"
var linkQuery = &models.RequestBodyTest{
	Url: "https://www.google.com/search?q=42&newwindow=1&sca_esv=585864680&source=hp&ei=2qllZfabGJevwPAPs92H4AE&iflsig=AO6bgOgAAAAAZWW36oesLy-IKv94a8R7Eg43XSviwif9&ved=0ahUKEwj28IKYqOaCAxWXFxAIHbPuARwQ4dUDCAo&uact=5&oq=42&gs_lp=Egdnd3Mtd2l6IgI0MjIFEC4YgAQyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIFEAAYgARI0RhQ3xRYzBZwAXgAkAEAmAFAoAFuqgEBMrgBA8gBAPgBAagCAA&sclient=gws-wiz",
}


func setupServer() *webserver.Router {
	dbStore := db.NewDBStore(testConfig)
	linksController := controllers.NewLinksController(kv.NewKVStore(testConfig), dbStore)
	recordsController := controllers.NewRecordsController(dbStore)

	return webserver.NewRouter(testConfig, linksController, recordsController)
}


func TestLinkCreation(t *testing.T) {
	var request *http.Request
	var w *httptest.ResponseRecorder
	var body []byte

	server := setupServer()

	// check empty db
	request, _ = http.NewRequest(http.MethodGet, "/" + alias, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	if w.Code != 404 {
		t.Errorf("error: should get 404 on this test, got %d", w.Code)
	}

	// create link
	body, _ = json.Marshal(linkQuery)
	request, _ = http.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	if w.Code != 200 {
		t.Errorf("error: should get 200 on this test, got %d", w.Code)
	}
	var resp models.ResponseBody
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	link := resp.Link

	// check redirect
	request, _ = http.NewRequest(http.MethodGet, "/" + link, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	if w.Code != 302 {
		t.Errorf("error: should get 302 on this test, got %d", w.Code)
	}
	redirectLocation, _ := w.Result().Location()
	if redirectLocation.String() != linkQuery.Url {
		t.Errorf("error: redirect url is wrong, is: %s, should be: %s", redirectLocation.String(), linkQuery.Url)
	}

	// create link for alias and check redirect
	var aliasQuery = &models.RequestBodyTest{
		Url: linkQuery.Url,
		Alias: alias,
	}
	body, _ = json.Marshal(aliasQuery)
	request, _ = http.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)

	request, _ = http.NewRequest(http.MethodGet, "/" + alias, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	if w.Code != 302 {
		t.Errorf("error: should get 302 on this test, got %d", w.Code)
	}
	redirectLocation, _ = w.Result().Location()
	if redirectLocation.String() != aliasQuery.Url {
		t.Errorf("error: redirect url is wrong, is: %s, should be: %s", redirectLocation.String(), aliasQuery.Url)
	}

	// attempt to create another link with the same alias
	body, _ = json.Marshal(aliasQuery)
	request, _ = http.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	if w.Code != 409 {
		t.Errorf("error: should get 409 on this test, got %d", w.Code)
	}

	// test getting records
	request, _ = http.NewRequest(http.MethodGet, "/records", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, request)
	var recs models.RecordsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &recs)

	cats := (*recs.Records[0]).Shorturl
	count := recs.Count
	if cats != "cats" || count != 2 {
		t.Errorf("error: link should be 'cats' and count should be 2")
	}
}