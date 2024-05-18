package main

import (
	"context"
	appConfig "kafka-server/internal/config"
	"kafka-server/internal/constants"
	"kafka-server/internal/consumer"
	database "kafka-server/internal/database"
	"kafka-server/internal/logging"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	logger, err := logging.NewLogger()

	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	appConfig, err := appConfig.NewConfig(logger)

	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	database_base_url := appConfig.Get("MONGO_URI")
	db_name := appConfig.Get("DB_NAME")
	main_server_collection_name := appConfig.Get("MAIN_SERVER_COLLECTION_NAME")
	cache_server_collection_name := appConfig.Get("CACHE_SERVER_COLLECTION_NAME")
	database_server_collection_name := appConfig.Get("DATABASE_SERVER_COLLECTION_NAME")

	topics := []string{constants.TOPIC_MAIN_SERVER,
		constants.TOPIC_CACHE_SERVER,
		constants.TOPIC_DATABASE_SERVER,
	}

	consumerGroup, err := sarama.NewConsumerGroup([]string{appConfig.Get("KAFKA_SERVICE_BASE_URL")}, "example-group", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	mongoMainServer, err := database.NewDbConnection(
		logger,
		database_base_url,
		db_name,
		main_server_collection_name,
	)

	if err != nil {
		log.Fatalf("Error creating mongo main server connection: %v", err)
	}

	mongoCacheServer, err := database.NewDbConnection(
		logger,
		database_base_url,
		db_name,
		cache_server_collection_name,
	)

	if err != nil {
		log.Fatalf("Error creating mongo cache server connection: %v", err)
	}

	mongoDatabaseServer, err := database.NewDbConnection(
		logger,
		database_base_url,
		db_name,
		database_server_collection_name,
	)

	if err != nil {
		log.Fatalf("Error creating mongo database server connection: %v", err)
	}

	consumer := consumer.NewConsumer(appConfig, mongoMainServer, mongoCacheServer, mongoDatabaseServer, logger)

	for {
		if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}
}
