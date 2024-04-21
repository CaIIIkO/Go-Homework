package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Event struct {
	Timestamp time.Time // время обращения
	Method    string    // название метода
	Request   string    // сырой запрос
}

type Kaf struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
}

var KafPrCo = Kaf{}

var (
	Broker = "localhost:9092"
	Topic  = "events"
)

func InitKafka() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	Producer, err := sarama.NewSyncProducer([]string{Broker}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}
	Consumer, err := sarama.NewConsumer([]string{Broker}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}
	KafPrCo = Kaf{
		Producer: Producer,
		Consumer: Consumer,
	}
}

func WriteToKafka(event Event, producer sarama.SyncProducer, topic string) error {
	eventBytes := []byte(fmt.Sprintf("%s %s: %s", event.Timestamp.Format(time.RFC3339), event.Method, event.Request))
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(eventBytes),
	})
	return err
}

func ReadFromKafka(consumer sarama.Consumer, topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %s\n", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error while receiving message: %v\n", err)
		}
	}
}
