package models

type Book struct {
	BookID		string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title 		string	`json:"title" bson:"title"`
	Author 		string 	`json:"author" bson:"author"`
	Publisher	string 	`json:"publisher" bson:"publisher"`
	Cost 		float32	`json:"cost" bson:"cost"`
}