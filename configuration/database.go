package configuration

import (
	"context"
	"go.uber.org/zap"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Database(url string, host string, port string, user string, passwd string, sch string, logger *zap.Logger) *pgxpool.Pool {
	url = strings.Replace(url, "{{host}}", host, -1)
	url = strings.Replace(url, "{{port}}", port, -1)
	url = strings.Replace(url, "{{username}}", user, -1)
	url = strings.Replace(url, "{{password}}", passwd, -1)
	url = strings.Replace(url, "{{schema}}", sch, -1)
	dbase, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		logger.Error("error when try to open database", zap.Error(err))
		panic(err)
	}
	return dbase
}
