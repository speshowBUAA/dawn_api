package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"fmt"
	"time"
	"go.uber.org/zap"
	"dawn_api/log"
)

//数据库地址
const (
    mongohost  = "localhost"
    mongoport  = 27017
	database   = "annotation_database"
	col        = "dawn"
)

type Student struct {
	Name string
}

var MongodbClient *mongo.Client
var MongodbCol *mongo.Collection

func NewMongoDBClient() {
	param := fmt.Sprintf("mongodb://%s:%d", mongohost, mongoport)
	clientOptions := options.Client().ApplyURI(param)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("Error", zap.Any("error", err))
	}

	log.Info("Connect to MongoDB!")
	MongodbClient = client
	collection := client.Database(database).Collection(col)
	MongodbCol = collection
}

func QueryMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cur, err := MongodbCol.Find(ctx, bson.D{})
	if err != nil {
		log.Error("Error", zap.Any("error", err))
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Error("Error", zap.Any("error", err))
		} 
		fmt.Println(result)
	}
}

func InsertMany(annos []interface{}) {
	insertResult, err := MongodbCol.InsertMany(context.TODO(), annos)
	if err != nil {
		log.Error("Error", zap.Any("error", err))
	}
	log.Info("Insert annos: ", zap.Any("ids", insertResult.InsertedIDs))
}

func CloseMongo() {
	err := MongodbClient.Disconnect(context.TODO())
	if err != nil {
		log.Error("Error", zap.Any("error", err))
	}
	log.Info("Connection to MongoDB closed.")
}