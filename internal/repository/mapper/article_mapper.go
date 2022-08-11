package mapper

import (
	"github.com/imamdnl/article/domain"
	"github.com/imamdnl/article/internal/repository/model"
)

func ToDomainArticle(m *model.ArticleModel) *domain.Article {
	return &domain.Article{
		Id:      m.Id,
		Author:  m.Author,
		Title:   m.Title,
		Body:    m.Body,
		Created: m.Created.Format("2006-01-02 15:04:05"),
	}
}

func ToDomainsArticle(m []*model.ArticleModel) []*domain.Article {
	var res []*domain.Article
	for i := range m {
		d := ToDomainArticle(m[i])
		res = append(res, d)
	}
	return res
}
