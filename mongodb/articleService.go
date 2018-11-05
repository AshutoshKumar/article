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

	for _, connect := range connects {
		index := getIndex(connect.Tags, tag)
		massagedData.Articles = append(massagedData.Articles, connect.ArticleId)
		connect.Tags = append(connect.Tags[:index], connect.Tags[index+1:]...)
		massagedData.Related_tags = RemoveDuplicatesFromSlice(append(massagedData.Related_tags, connect.Tags...))
	}
	return massagedData, err
}

func getIndex(tags []string, key string) int {
	for index, value := range tags {
		if value == key {
			return index
		}
	}
	return -1
}

func RemoveDuplicatesFromSlice(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
			fmt.Println(item, "is a duplicate")
		} else {
			m[item] = true
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}
