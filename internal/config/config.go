package config

import "os"

type Config struct {
	Env Env
}

type Env struct {
	PORT        string
	DB_URL      string
	BCRYPT_COST int `defualt:"10"`
	JWT_SECRET  string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Env: Env{
			PORT:        os.Getenv("PORT"),
			DB_URL:      os.Getenv("DB_URL"),
			BCRYPT_COST: 10,
			JWT_SECRET:  os.Getenv("JWT_SECRET"),
		},
	}

	return cfg, nil
}
