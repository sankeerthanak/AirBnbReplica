package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBName                 string
	Key                    string
	JWTExpirationInSeconds int64
	JWTSecret              string
	Region                 string
	S3BucketName           string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		DBUser:                 getEnv("DBUser", "mongoUser"),
		DBPassword:             getEnv("MongoDbPassword", "mongouser"),
		DBName:                 getEnv("DBName", "Airbnb"),
		Key:                    getEnv("KEY", "airbnbreplica"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 7*24*3600),
		JWTSecret:              getEnv("JWTSecret", "jwtsecret"),
		Region:                 getEnv("Region", "us-east-1"),
		S3BucketName:           getEnv("bucket", "airbnbreplica"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
