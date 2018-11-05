package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Body      string   `json:"body"`
	Tags      []string `json:"tags"`
	ArticleId string   `json:"-"`
}

type TagResponse struct {
	Tag          string   `json:"tag"`
	Count        int      `json:"count"`
	Articles     []string `json:"articles"`
	Related_tags []string `json:"related_tags"`
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
		Key:        []string{"articleid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// func ArticleModelTagsIndex() mgo.Index {
// 	return mgo.Index{
// 		Key:        []string{"tags"},
// 		Unique:     true,
// 		DropDups:   true,
// 		Background: true,
// 		Sparse:     true,
// 	}
// }

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
