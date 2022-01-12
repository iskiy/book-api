package main

import (
	"book-api/pkg/config"
	"book-api/pkg/handlers"
	"book-api/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

func main()  {
	config, err := config.NewConfig(".env")
	if err != nil{
		log.Fatalln(err)
	}

	mongo, err := repository.NewMongoDBClient(config.MongoConn, "book_api_db")
	if err != nil{
		log.Fatalln(err)
	}
	defer mongo.CloseConnection()

	rds := redis.NewClient(&redis.Options{
		Addr:		config.RedisConn,
		DB:          0,
	})

	h := handlers.NewHandlerManager(rds, mongo)

	r := gin.Default()
	bookAPI := r.Group("/api")
	{
		bookAPI.GET("/book/:id", h.GetBook)
		bookAPI.POST("/book", h.PutBook)
	}

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
