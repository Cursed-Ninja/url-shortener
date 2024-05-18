package consumer

import (
	"encoding/json"
	"errors"
	"kafka-server/internal/config"
	"kafka-server/internal/constants"
	"kafka-server/internal/database"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var (
	mongoMainServer     database.DBInterface
	mongoCacheServer    database.DBInterface
	mongoDatabaseServer database.DBInterface
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

type Job struct {
	message *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
}

type WorkerPool struct {
	jobs    chan Job
	results chan error
	wg      sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobs:    make(chan Job, numWorkers),
		results: make(chan error),
	}

	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (pool *WorkerPool) Submit(job Job) {
	pool.jobs <- job
}

func (pool *WorkerPool) worker() {
	defer pool.wg.Done()
	for job := range pool.jobs {
		var document interface{}
		var err error

		if err = json.Unmarshal(job.message.Value, &document); err != nil {
			job.session.MarkMessage(job.message, "")
			log.Printf("Error unmarshalling json: %v", err)
			pool.results <- err
			continue
		}

		v, _ := document.(map[string]interface{})
		v["expiresat"] = time.Now().AddDate(0, 0, 10)
		document = v

		switch job.message.Topic {
		case constants.TOPIC_CACHE_SERVER:
			err = mongoCacheServer.InsertOne(document)
		case constants.TOPIC_DATABASE_SERVER:
			err = mongoDatabaseServer.InsertOne(document)
		case constants.TOPIC_MAIN_SERVER:
			err = mongoMainServer.InsertOne(document)
		default:
			err = errors.New("invalid topic")
		}

		if err != nil {
			log.Printf("Error inserting document: %v", err)
			pool.results <- err
		} else {
			log.Printf("Message received: topic=%s partition=%d offset=%d key=%s value=%s\n",
				job.message.Topic, job.message.Partition, job.message.Offset, string(job.message.Key), string(job.message.Value))
			job.session.MarkMessage(job.message, "")
			pool.results <- nil
		}
	}
}

func (pool *WorkerPool) Shutdown() {
	close(pool.jobs)
	pool.wg.Wait()
	close(pool.results)
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
	mongoCacheServer = consumer.mongoCacheServer
	mongoDatabaseServer = consumer.mongoDatabaseServer
	mongoMainServer = consumer.mongoMainServer
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	err := mongoCacheServer.Disconnect()
	if err != nil {
		return err
	}

	err = mongoDatabaseServer.Disconnect()
	if err != nil {
		return err
	}

	err = mongoMainServer.Disconnect()
	if err != nil {
		return err
	}

	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	workerPool := NewWorkerPool(1000)

	go func() {
		for message := range claim.Messages() {
			workerPool.Submit(Job{message: message, session: session})
		}
		workerPool.Shutdown()
	}()

	for err := range workerPool.results {
		if err != nil {
			consumer.logger.Errorf("Processing error: %v", err)
		}
	}

	return nil
}
