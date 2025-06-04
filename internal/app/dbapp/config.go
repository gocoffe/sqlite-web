package dbapp

import (
	"github.com/antlko/golitedb/internal/db"
	"github.com/antlko/golitedb/internal/jwt"
	"github.com/antlko/golitedb/internal/server"
)

type Config struct {
	Hostname        string `env:"HOSTNAME"`
	ApplicationName string `env:"APPLICATION_NAME"`

	DB         db.Config
	Server     server.Config
	Authorizer jwt.Config
}
