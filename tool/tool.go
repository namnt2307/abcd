package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
)

var (
	KafkaProducer *kafka.Producer
)

type KafkaProducerModelStruct int

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var DataFile struct {
	Vod_ids []string
	Tag_ids []string
	Rib_ids []string
}

func ReadFile() {

	// Open our jsonFile
	jsonFile, err := os.Open("tool/data_init_seo_v5.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened data_init_seo_v5.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &DataFile)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	t := time.Now()
	ReadFile()

	// Init Kafka
	var serverKafka = "192.168.100.51:9092"
	var topicKafka = "TESTING3_MIGRATE_MYSQL_TO_MONGODB"
	var dataPushKafkaMigrate struct {
		Id       string
		Type     string
		Ref_type string
	}

	// Init Ribbon
	log.Println("++++++++++++++++++")
	log.Println("++++++++++++++++++")
	log.Println("Push Ribbon: Start")
	for k, v := range DataFile.Vod_ids {
		if k%10 == 0 {
			log.Println("Push Ribbon: ", k)
		}
		dataPushKafkaMigrate.Id = v
		dataPushKafkaMigrate.Type = "seo"
		dataPushKafkaMigrate.Ref_type = "vod"
		Send(topicKafka, serverKafka, dataPushKafkaMigrate)
		time.Sleep(time.Millisecond * 100)
	}
	log.Println("Push VOD: Done")

	latency := time.Since(t)
	log.Println("Duration (s):", latency.Seconds())
	time.Sleep(time.Second * 10)
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

func HandlePanic() {
	if err := recover(); err != nil {
		log.Println("HandlePanic:", err)
	}
	return
}
