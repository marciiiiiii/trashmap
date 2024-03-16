package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHelper struct {
	client *mongo.Client
}

func NewDatabaseHelper() *DatabaseHelper {
	return &DatabaseHelper{}
}

func (dbHelper *DatabaseHelper) Connect(databaseURI string) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal("panic in connect: ", err)
	}

	dbHelper.client = client

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
}

func (dbHelper *DatabaseHelper) FetchCollection(databaseName, collectionName string) ([]bson.M, error) {
	collection := dbHelper.client.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error Find: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
