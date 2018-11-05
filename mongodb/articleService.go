package mongo

import (
	"github.com/article/config"
	"github.com/article/model"
	"github.com/article/util"
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
	// collection.EnsureIndex(model.ArticleModelTagsIndex())
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
	model.Id = model.ArticleId
	return model, err
}

// GetArticleByTagName ...
func (p *ArticleService) GetArticleByTagName(tag string, date string) (model.TagResponse, error) {
	connects := []model.Article{}
	massagedData := model.TagResponse{}
	count, err := p.collection.Find(bson.M{"date": date}).Count()
	err = p.collection.Find(bson.M{"tags": tag}).All(&connects)
	massagedData.Tag = tag
	massagedData.Count = count

	//TODO we can move this to some util
	for _, connect := range connects {
		index := util.GetIndex(connect.Tags, tag)
		massagedData.Articles = append(massagedData.Articles, connect.ArticleId)
		connect.Tags = append(connect.Tags[:index], connect.Tags[index+1:]...)
		massagedData.Related_tags = util.RemoveDuplicatesFromSlice(append(massagedData.Related_tags, connect.Tags...))
	}
	return massagedData, err
}
