package handlers

import (
	. "book-api/pkg/models"
	. "book-api/pkg/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
)

var (
	redisWriteTime = 5 * time.Minute
)

type HandlerManager struct {
	db *MongoDBClient
	rds *redis.Client
}

func NewHandlerManager(rds *redis.Client, db *MongoDBClient) *HandlerManager {
	return &HandlerManager{rds: rds, db: db}
}

func (h *HandlerManager) GetBook(c *gin.Context){
	id := c.Param("id")

	rdsValue, err := h.rds.Get(id).Result()
	if !errors.Is(err, redis.Nil) {
		log.Println(fmt.Sprintf("Get from redis: %s", rdsValue))
		c.Data(http.StatusOK, "application/json", []byte(rdsValue))
		log.Println(err)
		return
	}

	book, err := h.db.GetBook(id)
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "get book by ID error")
	}

	jsonBook, err := json.Marshal(&book)
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "parse book error")
	}

	err = h.rds.Set(book.BookID, jsonBook, redisWriteTime).Err()
	if err != nil {
		log.Println(err.Error())
	}

	c.Data(http.StatusOK,"application/json", jsonBook)
}


func (h *HandlerManager) PutBook(c *gin.Context){
	book := Book{}
	err := c.Bind(&book)
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "get body error")
	}

	book, err = h.db.PutBook(book)
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "put book error")
	}

	jsonBook, err := json.Marshal(&book)
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "parse book error")
	}

	err = h.rds.Set(book.BookID, jsonBook, redisWriteTime).Err()
	if err != nil {
		log.Println(err.Error())
	}

	c.Data(http.StatusOK,"application/json", jsonBook)
}


func sendJSONError(c *gin.Context, status int, err error, message string){
	c.JSON(status, gin.H{
		"message": message,
	})
	log.Println(err)
	return
}

func (h *HandlerManager) DeleteDB() error{
	err := h.db.DeleteDB()
	if err != nil{
		return err
	}
	return h.db.CloseConnection()
}