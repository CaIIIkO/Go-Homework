package config

import (
	"github.com/redis/go-redis/v9"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "test"
	PASSWORD = "test"
	DBNAME   = "test"
)

type AuthConfigS struct {
	Username string
	Password string
}

var AuthConfig = AuthConfigS{
	Username: "user",
	Password: "user",
}

var RedisOpt = &redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}
