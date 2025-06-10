package dbapp

import (
	"github.com/gocoffe/sqlite-web/internal/db"
	"github.com/gocoffe/sqlite-web/internal/jwt"
	"github.com/gocoffe/sqlite-web/internal/server"
)

type Config struct {
	ApplicationName string `env:"APPLICATION_NAME"`

	DB         db.Config
	Server     server.Config
	Authorizer jwt.Config
}
