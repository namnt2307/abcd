package module

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
)

var (
	KafkaProducer *kafka.Producer
)

type KafkaProducerModelStruct int

func init() {
	var err error
	serverKafka, _ := CommonConfig.GetString("KAFKA", "uri")

	KafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": serverKafka})
	if err != nil {
		Sentry_log(err)
	}

	// defer KafkaProducer.Close()
	// Xử lý báo cáo gửi cho tin nhắn gửi đi
	go func() {
		defer HandlePanic()
		for e := range KafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					Sentry_log(err)
				}
			}
		}
	}()
}

func (this *KafkaProducerModelStruct) SendMessage(topic string, server string, message interface{}) bool {

	go func(topic string, server string, message interface{}) {
		defer HandlePanic()
		_, err := Send(topic, server, message)
		if err != nil {
			Sentry_log(err)
		}

	}(topic, server, message)

	return true
}

func Send(topic string, server string, message interface{}) (bool, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	dataByte, err := json.Marshal(message)
	if err != nil {
		return false, err
	}

	KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(string(dataByte)),
	}, nil)

	// Wait for message deliveries before shutting down
	KafkaProducer.Flush(15 * 1000)
	return true, nil
}
