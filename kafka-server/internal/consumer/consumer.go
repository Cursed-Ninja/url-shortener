package consumer

import (
	"encoding/json"
	"errors"
	"kafka-server/internal/config"
	"kafka-server/internal/constants"
	"kafka-server/internal/database"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerInterface interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type Consumer struct {
	mongoMainServer     database.DBInterface
	mongoCacheServer    database.DBInterface
	mongoDatabaseServer database.DBInterface
	config              config.ConfigInterface
	logger              *zap.SugaredLogger
}

func NewConsumer(config config.ConfigInterface, mongoMainServer database.DBInterface, mongoCacheServer database.DBInterface, mongoDatabaseServer database.DBInterface, logger *zap.SugaredLogger) *Consumer {
	return &Consumer{
		mongoMainServer:     mongoMainServer,
		mongoCacheServer:    mongoCacheServer,
		mongoDatabaseServer: mongoDatabaseServer,
		config:              config,
		logger:              logger,
	}
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var document interface{}
		var err error

		if err = json.Unmarshal([]byte(message.Value), &document); err != nil {
			consumer.logger.Errorf("Error unmarshalling json: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		switch message.Topic {
		case constants.TOPIC_CACHE_SERVER:
			err = consumer.mongoCacheServer.InsertOne(document)
		case constants.TOPIC_DATABASE_SERVER:
			err = consumer.mongoDatabaseServer.InsertOne(document)
		case constants.TOPIC_MAIN_SERVER:
			err = consumer.mongoMainServer.InsertOne(document)
		default:
			err = errors.New("invalid topic")
		}

		if err != nil {
			consumer.logger.Errorf("Error inserting document: %v", err)
		}

		consumer.logger.Infof("Message received: topic=%s partition=%d offset=%d key=%s value=%s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		session.MarkMessage(message, "")
	}
	return nil
}
