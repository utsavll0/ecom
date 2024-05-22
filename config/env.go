package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Config struct {
	DSN                    string
	JWTExpirationInSeconds int
	JWTSecret              string
}

var Envs = initConfig()

const projectDirName = "ecom"

func initConfig() Config {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
		DSN:                    getEnv("DSN", "postgres://utsav@localhost:5432/ecom?sslmode=disable"),
		JWTExpirationInSeconds: getEnvInt("JWTExpirationInSeconds", 3600*24*7),
		JWTSecret:              getEnv("JWTSecret", os.Getenv("JWTSecret")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		val, _ := strconv.Atoi(value)
		return val
	}
	return fallback
}
