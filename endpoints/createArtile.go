package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/article/model"
	mongo "github.com/article/mongodb"

	log "github.com/sirupsen/logrus"
)

// // Article ...
// type Articles struct {
// 	Id    int
// 	Title string
// 	Date  string
// 	Body  string
// 	Tags  []string
// }

func CreateArticleEndpoint(w http.ResponseWriter, r *http.Request, app *mongo.ArticleService) error {
	log.Infoln("[CreateArticle] >> [CreateArticleEndpoint] ")

	// if collection == nil {
	// 	return errors.New("no Mongodb connection available")
	// }

	article, err := validateArticleRequest(r)
	fmt.Print(article)
	if err != nil {
		log.Infoln("[CreateArticle] >> [CreateArticleEndpoint] Request validation failed" + err.Error())

		return err
		//Return Error from here
	}
	// err = app.CreateArticle(&article)
	return err
	//
	// err = mongo.CreateArticle(collection, &article)
	// if err != nil {
	// 	util.Error(w, http.StatusInternalServerError, err.Error())
	// 	return err
	// }
	//
	// return util.Json(w, "200", err)

}

func validateArticleRequest(r *http.Request) (model.ArticleModel, error) {
	log.Infoln("[CreateArticle] >> [validateArticleRequest] ")

	var article model.ArticleModel
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Failed to decode the body ")

		return model.ArticleModel{}, err
	}
	//No Error,lets validate all the fields

	// //Is Id blank ??
	// if article.Id {
	// 	log.Warnf("[CreateArticle] >>[validateArticleRequest] > Id is blank this is required field to create an Article")
	//
	// 	return nil, errors.New("[Article Id field is required to create an Article]")
	// }

	//Validate for Title, Data, Body fields
	if article.Title == "" || article.Date == "" || article.Body == "" {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Title|Data|Body fields are blank")

		return model.ArticleModel{}, errors.New("[Title|Date|body fields are blank]")
	}

	//Validation for tags field
	if len(article.Tags) == 0 {
		log.Warnf("[CreateArticle] >>[validateArticleRequest] > Tag field is blank ")

		return model.ArticleModel{}, errors.New("[Tag field is required to create an record")
	}

	return article, nil
}
