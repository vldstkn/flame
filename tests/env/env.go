package env

import (
	"database/sql"
	"flame/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"strconv"
)

type Env struct {
	*sql.DB
	Dsn        string
	ApiAddress string
	Jwt        string
}

func InitEnv() (*Env, error) {
	conf := config.LoadConfig("../../configs", "test")
	db, err := sql.Open("postgres", conf.Database.Dsn)
	if err != nil {
		return nil, err
	}
	return &Env{
		Dsn:        conf.Database.Dsn,
		DB:         db,
		Jwt:        conf.Auth.Jwt,
		ApiAddress: "http://" + conf.Public.Host + ":" + strconv.Itoa(conf.Public.Port) + "/api",
	}, nil
}

func (env *Env) Up() error {
	if err := goose.Up(env.DB, "../../migrations"); err != nil {
		return err
	}
	return nil
}
func (env *Env) Down() error {
	if err := goose.Down(env.DB, "../../migrations"); err != nil {
		return err
	}
	return nil
}
