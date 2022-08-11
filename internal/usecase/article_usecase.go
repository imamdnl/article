package usecase

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/fatih/structs"
	"github.com/imamdnl/article/dbmem"
	"github.com/imamdnl/article/domain"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type articleUseCase struct {
	repo   domain.ArticleRepository
	cache  dbmem.Cacher
	search *redisearch.Client
	logger *zap.Logger
}

func (a articleUseCase) ListArticle(request *domain.GetAllArticleRequest) interface{} {
	var out []*domain.Article
	var filAuthor string
	var filQry string
	if request.Author != "" {
		filAuthor = "@Author:" + request.Author + " "
	}
	if request.Query != "" {
		filQry = "(@Title:" + request.Query + ")|(@Body:" + request.Query + ")"
	}
	raw := filAuthor + filQry
	docs, _, _ := a.search.Search(
		redisearch.NewQuery(raw).
			SetSortBy("Created", false).
			Limit(0, 10).
			SetReturnFields("Id", "Author", "Title", "Body", "Created"))

	if len(docs) > 0 {
		for i := range docs {
			id, _ := strconv.Atoi(docs[i].Properties["Id"].(string))
			o := &domain.Article{
				Id:      id,
				Author:  docs[i].Properties["Author"].(string),
				Title:   docs[i].Properties["Title"].(string),
				Body:    docs[i].Properties["Body"].(string),
				Created: docs[i].Properties["Created"].(string),
			}
			out = append(out, o)
		}
		return out
	}

	res := a.repo.FindAllArticle(request)
	if res != nil {
		go func() {
			for _, v := range res {
				idConv := strconv.Itoa(v.Id)
				m := structs.Map(&v)
				a.cache.HMSet("article:"+idConv, idConv, m)
			}
		}()
	}
	return res
}

func (a articleUseCase) CreateArticle(r *domain.Article) interface{} {
	now := time.Now().Format("2006-01-02 15:04:05")
	r.Created = now
	res, err := a.repo.SaveArticle(r)
	if err == nil {
		r.Id = res
		idConv := strconv.Itoa(res)
		m := structs.Map(&r)
		go func() {
			a.cache.HMSet("article:"+idConv, idConv, m)
		}()
	}

	return res
}

func NewArticleUseCase(repo domain.ArticleRepository, cache dbmem.Cacher, search *redisearch.Client, logger *zap.Logger) domain.ArticleUseCase {
	return &articleUseCase{
		repo:   repo,
		cache:  cache,
		search: search,
		logger: logger,
	}
}
