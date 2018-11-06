package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/article/app"
	"github.com/article/config"
)

var validRequestArticlePayload = `{
    "id": "4",
    "title":"latest science shows that potato chips are better for you than sugar",
    "Date": "2016-09-22",
    "Body": "some text, potentially containing simple markup about how potato chips are great",
    "Tags":["health", "fitness", "science"]
}`

var invalidRequestArticlePayload = `{
    "id": "",
    "title":"",
    "Date": "2016-09-22",
    "Body": "some text, potentially containing simple markup about how potato chips are great",
    "Tags":["health", "fitness", "science"]
}`

var validArticleId = "4"
var invalidArticleId = "12345"
var validTagName = "health"
var validDate = "2016-09-22"
var invalidTagName = ""
var invalidDate = "2016-10-22"
var a app.App

func TestInitialize(t *testing.T) {
	Logger := log.New()
	log.SetFormatter(&log.JSONFormatter{})

	a.AppLogger = Logger

	a.Config = config.GetConfig()
	// var err error

	if a.Config.Mongo.Ip != "127.0.0.1:27017" {
		a.AppLogger.Fatalf("Enabled Mongodb deafult Ip does not match with expected value")
	}

	if a.Config.Mongo.DbName != "article" {
		a.AppLogger.Fatalf("Enabled Mongodb deafult dbname does not match with expected value")
	}
	if a.Config.Server.Port != ":8001" {
		a.AppLogger.Fatalf("Enabled Server Default port is not match with expected value")
	}

	a.Initialize()
}

func TestRequestPaths(t *testing.T) {

	t.Run("Test Create article path", func(t *testing.T) {
		t.Run("It should failed with validation error", func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/articles", ioutil.NopCloser(bytes.NewReader([]byte(invalidRequestArticlePayload))))
			req.Header.Set("Content-Type", "application/json")
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusBadRequest, response.Code)

		})

		t.Run("It should allow to create an article", func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/articles", ioutil.NopCloser(bytes.NewReader([]byte(validRequestArticlePayload))))
			req.Header.Set("Content-Type", "application/json")
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusOK, response.Code)

		})

		t.Run("It should failed with dubplicate record", func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/articles", ioutil.NopCloser(bytes.NewReader([]byte(validRequestArticlePayload))))
			req.Header.Set("Content-Type", "application/json")
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusBadRequest, response.Code)
		})
	})

	t.Run("Test getArticle based on passed article id", func(t *testing.T) {
		t.Run("It should return article", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/articles/"+validArticleId, nil)
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusOK, response.Code)

		})
		t.Run("It should failed when id is blank", func(t *testing.T) {

			req, _ := http.NewRequest("GET", "/articles/"+invalidArticleId, nil)
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusNotFound, response.Code)
		})
	})

	t.Run("Test getArticle based on passed tag", func(t *testing.T) {
		t.Run("It should return list of article", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/tags/"+validTagName+"/"+validDate, nil)
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusOK, response.Code)
		})
		t.Run("It should failed with bad request ", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/tags/"+invalidTagName+"/"+invalidDate, nil)
			response := executeRequest(a.Router, req)
			checkResponseCode(t, http.StatusBadRequest, response.Code)
		})
	})
}

func executeRequest(router *chi.Mux, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected Status Code:%d | Received Status Code:%d", expected, actual)
	}
}
