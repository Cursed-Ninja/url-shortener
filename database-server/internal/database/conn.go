package database

import (
	"context"
	"url-shortner-database/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DBInterface interface {
	InsertOne(document models.URL) error
	FindOne(filter bson.D) (models.URL, error)
}

type dB struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

func NewDbConnection(logger *zap.SugaredLogger, dbConnection string, dbName string, collectionName string) (*dB, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbConnection).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		logger.Errorw("Could not connect to mongo", zap.Error(err))
		return nil, err
	}

	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		logger.Errorw("Could not ping database", zap.Error(err))
		return nil, err
	}

	indexOptions := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "expiresat", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
		{
			Keys:    bson.D{{Key: "shorturlpath", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	collection := client.Database(dbName).Collection(collectionName)

	collection.Indexes().CreateMany(context.TODO(), indexOptions)

	logger.Infow("Successfully established connection")

	return &dB{
		collection: collection,
		logger:     logger,
		client:     client,
	}, nil
}

func (connection *dB) InsertOne(document models.URL) error {
	_, err := connection.collection.InsertOne(context.TODO(), document)

	if err != nil {
		connection.logger.Errorw("Could not insert document", zap.Error(err), zap.Any("document", document))
	}

	return err
}

func (connection *dB) FindOne(filter bson.D) (models.URL, error) {
	var result models.URL
	err := connection.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != mongo.ErrNoDocuments {
		connection.logger.Errorw("Error retrieving documents", zap.Error(err))
	}

	return result, err
}

func (connection *dB) Disconnect() error {
	err := connection.client.Disconnect(context.TODO())

	if err != nil {
		connection.logger.Errorw("Could not disconnect from database", zap.Error(err))
		return err
	}

	return nil
}

func (connection *dB) DeleteDb(databaseName string) error {
	err := connection.client.Database(databaseName).Drop(context.Background())

	if err != nil {
		connection.logger.Errorw("Could not drop database", zap.Error(err))
		return err
	}

	return nil
}
