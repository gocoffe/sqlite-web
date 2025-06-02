package server

type Config struct {
	Port string `env:"SERVER_PORT,default=5577"`
}
