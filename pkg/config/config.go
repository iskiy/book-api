package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)
type Config struct {
	RedisConn	string
	MongoConn	string
}

func NewConfig(configSource string) (*Config, error) {
	err := godotenv.Load(configSource)
	if err != nil{
		return nil, err
	}

	redisConn, ok := os.LookupEnv("REDIS_CONN")
	if !ok{
		return nil, errors.New("can`t find REDIS_CONN env variable")
	}

	mongoConn, ok := os.LookupEnv("MONGO_CONN")
	if !ok{
		return nil, errors.New("can`t find MONGO_CONN env variable")
	}

	return &Config{RedisConn: redisConn, MongoConn: mongoConn}, nil
}
