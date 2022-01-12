package repository

import (
	. "book-api/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Repository interface {
	GetBook(bookID string) (Book, error)
	PutBook(book Book) (*Book, error)
}


type MongoDBClient struct {
	booksCollection *mongo.Collection
}
func NewMongoDBClient(connString, dbName string) (*MongoDBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil{
		return nil, err
	}

	ctx := context.TODO()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	apiDB := client.Database(dbName)
	booksCollection := apiDB.Collection("books")

	return &MongoDBClient{ booksCollection: booksCollection}, nil
}

func (db *MongoDBClient) PutBook(book Book) (Book,error){
	ctx := context.TODO()

	resultID, err := db.booksCollection.InsertOne(ctx, book)
	if err != nil {
		return Book{}, err
	}

	book.BookID = resultID.InsertedID.(primitive.ObjectID).Hex()
	return book, nil
}


func (db *MongoDBClient) GetBook(bookID string) (*Book,error){
	id, err := primitive.ObjectIDFromHex(bookID)
	if err != nil{
		return nil, err
	}

	var book Book
	ctx := context.TODO()
	err = db.booksCollection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (db *MongoDBClient) CloseConnection() error {
	return db.booksCollection.Database().Client().Disconnect(context.TODO())
}

func (db *MongoDBClient) DeleteDB() error{
	err := db.booksCollection.Database().Drop(context.TODO())
	if err != nil{
		return err
	}
	return nil
}