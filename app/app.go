package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/article/config"
	"github.com/article/model"
	mongo "github.com/article/mongodb"
	"github.com/article/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

//App implements the core pieces of the api, including integrations between functions, models and routes
type App struct {
	Router         *chi.Mux // Handles server routes
	session        *mongo.Session
	AppLogger      *log.Logger
	Config         *config.Config
	articleService *mongo.ArticleService
}

//Create Articles
func (a *App) createArticle(w http.ResponseWriter, r *http.Request) {
	// initalize variables to check throughout request
	a.AppLogger.Infoln("In create Article")

	article, err := validateArticleRequest(r)
	if err != nil {
		log.Infoln("[CreateArticle] >> [CreateArticleEndpoint] Request validation failed" + err.Error())

		//Return Error from here
		util.Error(w, http.StatusInternalServerError, err.Error())
	}

	articleRes := a.articleService.CreateArticle(&article)

	if articleRes != nil {
		util.Error(w, http.StatusInternalServerError, articleRes.Error())
		return
	}

	util.Json(w, http.StatusOK, articleRes)
}

func validateArticleRequest(r *http.Request) (model.Article, error) {
	log.Infoln("[CreateArticle] >> [validateArticleRequest] ")

	var article model.Article
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Failed to decode the body ")

		return model.Article{}, err
	}
	//No Error,lets validate all the fields

	//Is Id blank ??
	if article.Id == "" {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Id is blank this is required field to create an Article")

		return model.Article{}, errors.New("[Article Id field is required to create an Article]")
	}

	//Validate for Title, Data, Body fields
	if article.Title == "" || article.Date == "" || article.Body == "" {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Title|Data|Body fields are blank")

		return model.Article{}, errors.New("[Title|Date|body fields are blank]")
	}

	//Validation for tags field
	if len(article.Tags) == 0 {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Tag field is blank ")

		return model.Article{}, errors.New("[Tag field is required to create an record")
	}

	return article, nil
}

func (a *App) getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "id")
	articles, err := a.articleService.GetArticleById(articleID)

	if err != nil {
		util.Error(w, http.StatusNotFound, err.Error())
		return
	}

	util.Json(w, http.StatusOK, articles)
}

func (a *App) getTagHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("REached here")

	tagName := chi.URLParam(r, "tagName")
	date := chi.URLParam(r, "date")
	articles, err := a.articleService.GetArticleByTagName(tagName, date)

	if err != nil {
		util.Error(w, http.StatusNotFound, err.Error())
		return
	}

	util.Json(w, http.StatusOK, articles)
}

//Initialize bootstaps database, preps cache and creates initial router
func (a *App) Initialize() {
	Logger := log.New()
	log.SetFormatter(&log.JSONFormatter{})
	a.AppLogger = Logger

	a.Config = config.GetConfig()
	var err error

	a.session, err = mongo.NewSession(a.Config.Mongo)

	if err != nil {
		a.AppLogger.Fatalf("Error connecting to mongodb. Following Error passed: %v", err)
	}

	a.articleService = mongo.NewArticleService(a.session.Copy(), a.Config.Mongo)

	if err != nil {
		a.AppLogger.Fatalf("Error connecting to collections. Following Error passed: %v", err)
	}

	// create router and add routes
	a.AppLogger.Println("creating router...")
	a.Router = chi.NewRouter()
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.AppLogger.Println("router created!")
	a.InitializeRoutes()
}

//InitializeRoutes assigns handlers for app methods
func (a *App) InitializeRoutes() {
	a.AppLogger.Infoln("adding routes...")

	a.Router.HandleFunc("/articles", a.createArticle)
	a.Router.MethodFunc("GET", "/getarticles/{id}", a.getArticlesHandler)
	a.Router.MethodFunc("GET", "/tags/{tagName}/{date}", a.getTagHandler)
	a.AppLogger.Infoln("routes added!")
}

// Run ...
func (a *App) Run() {
	fmt.Println("Run")
	defer a.session.Close()
	a.Start()
}

// Start ...
func (a *App) Start() {
	// start server
	a.RunApp(a.Config.Server.Port)
}

//RunApp ...
func (a *App) RunApp(add string) {
	a.AppLogger.Infoln("starting web server")
	// create and run basic http server
	httpServer := &http.Server{
		Addr:              add,
		Handler:           a.Router,
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
	}
	a.AppLogger.Fatalln(httpServer.ListenAndServe())
	a.AppLogger.Infoln("web server running")
}
