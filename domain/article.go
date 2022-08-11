package domain

type Article struct {
	Id      int    `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type GetAllArticleRequest struct {
	Query  string `json:"query"`
	Author string `json:"author"`
}

type ArticleRepository interface {
	FindAllArticle(request *GetAllArticleRequest) []*Article
	SaveArticle(d *Article) (int, error)
}

type ArticleUseCase interface {
	ListArticle(request *GetAllArticleRequest) interface{}
	CreateArticle(r *Article) interface{}
}
