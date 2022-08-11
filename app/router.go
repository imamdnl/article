package app

import (
	ctrl "github.com/imamdnl/article/internal/delivery/http"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(articleHttp ctrl.ArticleHttp) *httprouter.Router {
	router := httprouter.New()
	router.POST("/articles", articleHttp.CreateArticle)
	router.GET("/articles", articleHttp.FindAllArticle)

	return router
}
