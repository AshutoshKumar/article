package mongo

import (
	"fmt"

	"github.com/article/config"
	"github.com/article/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ArticleService struct {
	collection *mgo.Collection
}

// NewArticleService ...
func NewArticleService(session *mgo.Session, config *config.MongoConfig) *ArticleService {
	collection := session.DB(config.DbName).C("article")
	collection.EnsureIndex(model.ArticleModelIndex())
	return &ArticleService{collection}
}

// CreateArticle ...
func (p *ArticleService) CreateArticle(u *model.Article) error {
	article, err := model.NewArticleModel(u)
	if err != nil {
		return err
	}
	return p.collection.Insert(&article)
}

// GetArticleById ...
func (p *ArticleService) GetArticleById(id string) (model.Article, error) {
	model := model.Article{}
	err := p.collection.Find(bson.M{"articleid": id}).One(&model)
	fmt.Print(model)
	fmt.Print(err)
	return model, err
}

// GetArticleByTagName ...
func (p *ArticleService) GetArticleByTagName(tag string, date string) (model.Article, error) {
	model := model.Article{}
	err := p.collection.Find(bson.M{"tag": tag}).One(&model)
	return model, err
}
