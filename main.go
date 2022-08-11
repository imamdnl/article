package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imamdnl/article/app"
	"github.com/imamdnl/article/configuration"
	"github.com/imamdnl/article/dbmem"
	"github.com/imamdnl/article/domain"
	ctrl "github.com/imamdnl/article/internal/delivery/http"
	"github.com/imamdnl/article/internal/repository"
	"github.com/imamdnl/article/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	logger := configuration.Logger()
	configuration.Environment(logger)

	var adrs []string
	rdp := os.Getenv("REDIS_PORT")
	rdpwd := os.Getenv("REDIS_PASSWORD")
	for _, adr := range strings.Split(os.Getenv("REDIS_HOST"), ",") {
		adrs = append(adrs, adr+":"+rdp)
	}
	cutil := dbmem.NewCache(
		configuration.ConfigCache(adrs[0], rdpwd),
		configuration.ConfigClusterCache(adrs, rdpwd),
		os.Getenv("REDIS_CLUSTER") == "TRUE", logger)

	redisearch := configuration.ConfigSearch(adrs[0])

	db := configuration.Database(
		os.Getenv("DB_URL"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		logger,
	)
	repoBase := domain.NewBaseRepository(db)

	repoArticle := repository.NewArticleRepository(*repoBase, logger)
	article := usecase.NewArticleUseCase(repoArticle, cutil, redisearch, logger)
	articleCtrl := ctrl.NewArticleHttp(article, logger)

	router := app.NewRouter(articleCtrl)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: router,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Error("exit reason: ", zap.Error(err))
	}
}

func NewHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	d := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    200,
		Message: "OK",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		panic(err)
	}
}
