package database_test

import (
	"encoding/json"
	"kafka-server/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
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

	jsonString := `{"level":"error","ts":1716041809.434134,"msg":"{requestId 15 0 7468b69a-28f4-4dca-af5d-52a468475e28 <nil>}Error sending request to database service{error 26 0  Post \"http://localhost:8081/redirect\": dial tcp [::1]:8081: connectex: No connection could be made because the target machine actively refused it.}"}`

	var document interface{}

	if err = json.Unmarshal([]byte(jsonString), &document); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
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
