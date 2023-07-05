package config

import "os"

type EnvConfig struct {
	Port       string
	DbHost     string
	DbPort     string
	DbUser     string
	DbName     string
	DbPassword string
	DbSchema   string
}

func SetEnv() {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "test2")
	os.Setenv("DB_PASSWORD", "2002")
	os.Setenv("DB_SCHEMA", "disable")
}

func GetEnv() *EnvConfig {
	return &EnvConfig{
		Port:       os.Getenv("PORT"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbName:     os.Getenv("DB_NAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbSchema:   os.Getenv("DB_SCHEMA"),
	}
}
