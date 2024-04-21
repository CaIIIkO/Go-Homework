package config

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
