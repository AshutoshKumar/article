package endpoints

import (
	"errors"
	"net/http"

	"github.com/article/model"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetArticles ...
func GetArticle(w http.ResponseWriter, r *http.Request, session *mgo.Session) (model.ArticleModel, error) {
	log.Infoln("[getARticles] >> [getArticles] ")

	if session == nil {
		return model.ArticleModel{}, errors.New("no Mongodb connection available")
	}

	articleID, err := WellFormedKey(r.URL.Query().Get("id"), r.Header.Get("id"))

	if err != nil {
		return model.ArticleModel{}, err
	}

	articleRes, err := FindArticleById(articleID, session)

	if err != nil {
		return model.ArticleModel{}, err
	}

	return articleRes, nil
	//Extarct Articles ID from the request parameter
}

func WellFormedKey(queryKey string, headerKey string) (string, error) {
	hasQuery := queryKey != ""
	hasHeader := headerKey != ""

	if !hasQuery && !hasHeader {
		return "", errors.New("No key found in request")
	}

	if hasQuery && hasHeader && queryKey != headerKey {
		return "", errors.New("Article id in header and query do not match")
	}
	return queryKey, nil

}

func FindArticleById(id string, session *mgo.Session) (model.ArticleModel, error) {
	var article model.ArticleModel
	err := session.DB("article").C("articles").FindId(bson.M{"id": id}).One(&article)
	if err != nil {
		return model.ArticleModel{}, err
	}

	return article, nil
}
