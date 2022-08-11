package http

import (
	"encoding/json"
	"github.com/imamdnl/article/domain"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

type ArticleHttp interface {
	CreateArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	FindAllArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type articleHttp struct {
	articleUsecase domain.ArticleUseCase
	logger         *zap.Logger
}

func NewArticleHttp(articleUseCase domain.ArticleUseCase, logger *zap.Logger) ArticleHttp {
	return &articleHttp{
		articleUsecase: articleUseCase,
		logger:         logger,
	}
}

func (s articleHttp) CreateArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req domain.Article
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}

	articleResponse := s.articleUsecase.CreateArticle(&req)
	webResponse := domain.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   articleResponse,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(webResponse)
	if err != nil {
		panic(err)
	}
}

func (s articleHttp) FindAllArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var req domain.GetAllArticleRequest
	req.Author = r.URL.Query().Get("author")
	req.Query = r.URL.Query().Get("query")

	articleResponse := s.articleUsecase.ListArticle(&req)
	webResponse := domain.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   articleResponse,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(webResponse)
	if err != nil {
		panic(err)
	}
}
