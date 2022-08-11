package repository

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/imamdnl/article/domain"
	"github.com/imamdnl/article/internal/repository/mapper"
	"github.com/imamdnl/article/internal/repository/model"
	"go.uber.org/zap"
)

type articleRepository struct {
	super  domain.BaseCapsule
	logger *zap.Logger
}

func (a articleRepository) FindAllArticle(r *domain.GetAllArticleRequest) (res []*domain.Article) {
	var out []*model.ArticleModel
	var and1 string
	var and2 string
	if r.Author != "" {
		and1 = " and author like '%" + r.Author + "%' "
	}
	if r.Query != "" {
		and2 = " and (title like '%" + r.Query + "%' or body like '%" + r.Query + "%') "
	}
	query := "select id, author, title, body, created from article where 1=1 " + and1 + and2 + " order by created desc limit 10 "
	err := pgxscan.Select(context.Background(), a.super.Database, &out, query)

	if err != nil {
		a.logger.Error("Article get all data failed with error = ", zap.Error(err))
		return nil
	}
	return mapper.ToDomainsArticle(out)
}

func (a articleRepository) SaveArticle(d *domain.Article) (int, error) {
	var id int
	tx, txErr := a.super.Database.Begin(context.Background())
	if txErr != nil {
		a.logger.Error("Begin transaction to save article is failed = ", zap.Error(txErr))
		return id, fmt.Errorf("begin transaction to save article is failed = %s", txErr.Error())
	}
	defer tx.Rollback(context.TODO())

	articleId := tx.QueryRow(context.Background(), "insert into article (author, title, body, created) values ($1,$2,$3,$4) RETURNING ID", d.Author, d.Title, d.Body, d.Created)

	err := articleId.Scan(&id)
	if err != nil {
		a.logger.Error("Create Article failed with error = ", zap.Error(err))
		return id, fmt.Errorf("create Article failed with error = %s", txErr.Error())
	}

	if txErr = tx.Commit(context.Background()); txErr != nil {
		a.logger.Error("failed to create a article with error = ", zap.Error(txErr))
		return id, fmt.Errorf("failed to create a article with error = %s", txErr.Error())
	}
	return id, nil
}

func NewArticleRepository(super domain.BaseCapsule, logger *zap.Logger) domain.ArticleRepository {
	return &articleRepository{
		super:  super,
		logger: logger,
	}
}
