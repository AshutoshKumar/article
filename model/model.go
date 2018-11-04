package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Id        string
	Title     string
	Date      string
	Body      string
	Tags      []string
	ArticleId string
}

// ArticleModel ...
type ArticleModel struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	ArticleId string
	Title     string
	Body      string
	Tags      []string
	Date      string
}

// ArticleModelIndex ...
func ArticleModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"ArticleId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// NewArticleModel ...
func NewArticleModel(u *Article) (*ArticleModel, error) {
	article := ArticleModel{
		ArticleId: u.Id,
		Title:     u.Title,
		Body:      u.Body,
		Tags:      u.Tags,
		Date:      u.Date,
	}

	return &article, nil
}
