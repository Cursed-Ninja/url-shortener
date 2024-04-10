package database_test

import (
	"testing"
	"time"
	"url-shortner-database/internal/database"
	"url-shortner-database/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

var testStruct struct {
	logger           *zap.SugaredLogger
	connectionString string
	connectionDb     string
	connectionColl   string
}

func TestMain(m *testing.M) {
	testStruct.logger = zap.NewNop().Sugar()
	testStruct.connectionString = "" // Enter a valid connection string
	testStruct.connectionDb = "testDb"
	testStruct.connectionColl = "test"
	m.Run()
}

func TestNewDbConnection(t *testing.T) {
	t.Run("Invalid Case", func(t *testing.T) {
		_, err := database.NewDbConnection(testStruct.logger, "invalid", testStruct.connectionDb, testStruct.connectionColl)
		assert.NotNil(t, err, "Error creating db connection")
	})

	t.Run("Valid Case", func(t *testing.T) {
		db, err := database.NewDbConnection(testStruct.logger, testStruct.connectionString, testStruct.connectionDb, testStruct.connectionColl)
		assert.Nil(t, err, "Error creating db connection")
		db.DeleteDb(testStruct.connectionDb)
		db.Disconnect()
	})
}

func TestInsertOne(t *testing.T) {
	db, err := database.NewDbConnection(testStruct.logger, testStruct.connectionString, testStruct.connectionDb, testStruct.connectionColl)

	if err != nil {
		t.Fatalf("Error creating db connection: %v", err)
	}

	document := models.URL{
		ShortUrlPath: "test",
		OriginalUrl:  "test",
		ExpiresAt:    time.Now(),
	}

	err = db.InsertOne(document)
	assert.Nil(t, err, "Error inserting document")

	err = db.DeleteDb(testStruct.connectionDb)

	if err != nil {
		t.Fatalf("Error deleting database: %v", err)
	}

	err = db.Disconnect()

	if err != nil {
		t.Fatalf("Error disconnecting: %v", err)
	}
}

func TestFindOne(t *testing.T) {
	db, err := database.NewDbConnection(testStruct.logger, testStruct.connectionString, testStruct.connectionDb, testStruct.connectionColl)

	if err != nil {
		t.Fatalf("Error creating db connection: %v", err)
	}

	t.Run("Not found case", func(t *testing.T) {
		filter := bson.D{{Key: "shortenedurl", Value: "testDb"}}
		_, err := db.FindOne(filter)
		assert.NotNil(t, err, "Error finding document")
	})

	document := models.URL{
		ShortUrlPath: "test",
		OriginalUrl:  "test",
		ExpiresAt:    time.Now(),
	}

	err = db.InsertOne(document)

	if err != nil {
		t.Fatalf("Error inserting document: %v", err)
	}

	t.Run("Found case", func(t *testing.T) {
		filter := bson.D{{Key: "shortenedurl", Value: "test"}}
		_, err := db.FindOne(filter)
		assert.Nil(t, err, "Error finding document")
	})

	err = db.DeleteDb(testStruct.connectionDb)

	if err != nil {
		t.Fatalf("Error deleting database: %v", err)
	}

	err = db.Disconnect()

	if err != nil {
		t.Fatalf("Error disconnecting: %v", err)
	}
}
