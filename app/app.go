package app

import (
	"encoding/json"
	"errors"
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
	Session        *mongo.Session
	AppLogger      *log.Logger
	Config         *config.Config
	articleService *mongo.ArticleService
}

//Initialize bootstaps database, preps cache and creates initial router
func (a *App) Initialize() {
	Logger := log.New()
	log.SetFormatter(&log.JSONFormatter{})
	a.AppLogger = Logger

	a.Config = config.GetConfig()
	var err error

	a.Session, err = mongo.NewSession(a.Config.Mongo)

	if err != nil {
		a.AppLogger.Fatalf("Error connecting to mongodb. Following Error passed: %v", err)
	}

	a.articleService = mongo.NewArticleService(a.Session.Copy(), a.Config.Mongo)

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
	a.AppLogger.Infoln("[app] >> [InitializeRoutes] > adding routes...")

	a.Router.HandleFunc("/articles", a.createArticle)
	a.Router.MethodFunc("GET", "/articles/{id}", a.getArticlesHandler)
	a.Router.MethodFunc("GET", "/tags/{tagName}/{date}", a.getTagHandler)
	a.AppLogger.Infoln("routes added!")
}

// Run ...
func (a *App) Run() {
	a.AppLogger.Debugln("[app] >> [Run]")
	defer a.Session.Close()
	a.Start()
}

// Start ...
func (a *App) Start() {
	a.AppLogger.Debugln("[app] >> [Start]")

	// start server
	a.RunApp(a.Config.Server.Port)
}

//RunApp ...
func (a *App) RunApp(add string) {
	a.AppLogger.Infoln("[app] >> [RunApp] > starting web server")
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

// This function we be get used to create an article
func (a *App) createArticle(w http.ResponseWriter, r *http.Request) {
	// initalize variables to check throughout request
	a.AppLogger.Debugln("[app] >>[createArticle]")

	article, err := validateArticleRequest(r)
	if err != nil {
		a.AppLogger.Warnln("[app] >>[createArticle] > Failed to validate request")

		//Return Error from here
		util.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	articleRes := a.articleService.CreateArticle(&article)

	if articleRes != nil {
		a.AppLogger.Warnln("[app] >>[CreateArticle] > Failed to create an article" + articleRes.Error())

		util.Error(w, http.StatusBadRequest, articleRes.Error())
		return
	}
	a.AppLogger.Infoln("Article created Successfully")

	util.Json(w, http.StatusOK, "created")
}

//This function will be get used to validate the request body for create an article
func validateArticleRequest(r *http.Request) (model.Article, error) {
	log.Infoln("[app] >> [validateArticleRequest]")

	var article model.Article
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		log.Warnf("[app] >>[validateArticleRequest] > Failed to decode the body ")

		return model.Article{}, err
	}

	//No Error,lets validate all the fields

	//Is Id blank ??
	if article.Id == "" {
		log.Warnf("[app] >>[validateArticleRequest] > Id is blank this is required field to create an Article")

		return model.Article{}, errors.New("[Article Id field is required to create an Article]")
	}

	//Validate for Title, Data, Body fields
	if article.Title == "" || article.Date == "" || article.Body == "" {
		log.Warnf("[app] >>[validateArticleRequest] > Title|Data|Body fields are blank")

		return model.Article{}, errors.New("[Title|Date|body fields are blank]")
	}

	//Todo we should handle date check properly like @Ashutosh
	//Validation for tags field
	if len(article.Tags) == 0 {
		log.Warnf("[app] >>[validateArticleRequest] > Tag field is blank ")

		return model.Article{}, errors.New("[Tag field is required to create an record")
	}

	return article, nil
}

//This fucntion will be get used to get Article based on the passed article id
func (a *App) getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	a.AppLogger.Debugln("[app] >>[getArticlesHandler]")

	//get id from the request url
	articleID := chi.URLParam(r, "id")

	//Is articleId available ?
	if articleID == "" {
		a.AppLogger.Warnln("[app] >>[getArticlesHandler] > article id is blank")

		util.Error(w, http.StatusBadRequest, "id is missing")
		return
	}

	articles, err := a.articleService.GetArticleById(articleID)

	if err != nil {
		a.AppLogger.Warnln("[app] >>[getArticlesHandler] > failed to fetch article" + err.Error())

		util.Error(w, http.StatusNotFound, err.Error())
		return
	}
	a.AppLogger.Infoln("Article available for passed article id" + articleID)

	util.Json(w, http.StatusOK, articles)
}

func (a *App) getTagHandler(w http.ResponseWriter, r *http.Request) {
	a.AppLogger.Debugln("[app] >>[getTagHandler]")

	//Extract TagName and date from the request url
	tagName := chi.URLParam(r, "tagName")
	date := chi.URLParam(r, "date")

	//Tagname and data ??
	if tagName == "" || date == "" {
		a.AppLogger.Warnln("[app] >>[getTagHandler] > tagname|date field are blank")

		util.Error(w, http.StatusBadRequest, "tagname or date are blank")
		return
	}

	date = date[0:4] + "-" + date[4:6] + "-" + date[6:]
	articles, err := a.articleService.GetArticleByTagName(tagName, date)

	if err != nil {
		a.AppLogger.Warnln("[app] >>[getTagHandler] > failed to fetch article based on passed tag and date field")

		util.Error(w, http.StatusNotFound, err.Error())
		return
	}

	a.AppLogger.Infoln("Article available for passed tagname and date" + "tagname" + tagName + "date" + date)
	util.Json(w, http.StatusOK, articles)
}
