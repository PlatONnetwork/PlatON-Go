package platonstats

import (
	"os"
	"testing"
	"time"

	cluster "github.com/bsm/sarama-cluster"

	"github.com/Shopify/sarama"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

func consumerSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_5_0_0
	//手工提交offset
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	//初始从最新的offset开始
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}

func Test_kafkaClient_producer1(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	blockProducer, err := sarama.NewSyncProducer([]string{"192.168.112.32:9092"}, producerConfig())
	if err != nil {
		log.Error("Failed to init msg Kafka sync producer....", "err", err)
	}

	//kafkaClient := NewKafkaClient("192.168.112.32:9092", "", "")

	log.Info("Success to init msg Kafka client ....")

	msg := &sarama.ProducerMessage{
		Topic:     "test-sender2",
		Partition: 0,
		Key:       sarama.StringEncoder("key-1"),
		Value:     sarama.StringEncoder("value-1"),
		Timestamp: time.Now(),
	}

	partition, offset, err := blockProducer.SendMessage(msg)

	if err != nil {
		log.Error("Failed to send message to Kafka", "error", err)
	} else {
		log.Info("Success to send message to Kafka", "partition", partition, "offset", offset)
	}
}

func Test_kafkaClient_producer2(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	produceTopic := "test-sender4"
	kafkaClient := NewKafkaClient("192.168.112.32:9092", produceTopic, "")

	log.Info("Success to init msg Kafka client ....")

	msg := &sarama.ProducerMessage{
		Topic:     produceTopic,
		Partition: 0,
		Key:       sarama.StringEncoder("key-1"),
		Value:     sarama.StringEncoder("value-1"),
		Timestamp: time.Now(),
	}

	partition, offset, err := kafkaClient.syncProducer.SendMessage(msg)

	if err != nil {
		log.Error("Failed to send message to Kafka", "error", err)
	} else {
		log.Info("Success to send message to Kafka", "partition", partition, "offset", offset)
	}
}

func Test_kafkaClient_consumer(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	kafkaClient := NewKafkaClient("192.168.112.32:9092", "", "")

	log.Info("Success to init msg Kafka account-checking consumer....")

	for {
		select {
		case msg := <-kafkaClient.partitionConsumer.Messages():
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message", "offset", msg.Offset, "key", key, "value", value)

			//手工提交offset()
			kafkaClient.partitionOffsetManager.MarkOffset(msg.Offset, "")
		case err := <-kafkaClient.partitionConsumer.Errors():
			log.Error("Failed to pull account-checking message from Kafka", "err", err)
			panic(err)
		}
	}
}

func T2st_kafkaMsg_work(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	brokers := []string{"192.168.112.32:9092"}
	topic := "platon-account-checking"
	if len(topic) == 0 {
		topic = defaultKafkaAccountCheckingTopic
	}

	client, err := sarama.NewClient(brokers, saramaConfig())
	if err != nil {
		panic("client create error")
	}
	defer client.Close()
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		panic("offsetManager create error")
	}
	defer consumer.Close()

	offsetManager, err := sarama.NewOffsetManagerFromClient("group1", client)
	if err != nil {
		panic("offsetManager create error")
	}
	defer offsetManager.Close()
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, int32(0))
	if err != nil {
		panic("offsetManager create error")
	}
	defer partitionOffsetManager.Close()

	checkingConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Error("Failed to create Kafka partition_consumer....", "err", err)
		panic(err)
	}
	defer func() {
		if err := checkingConsumer.Close(); err != nil {
			log.Error("Failed to close checkingConsumer", "err", err)
		}
	}()

	log.Info("Success to init msg Kafka account-checking consumer....")
	for {
		select {
		case msg := <-checkingConsumer.Messages():
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message", "offset", msg.Offset, "key", key, "value", value)

			//手工提交offset()
			partitionOffsetManager.MarkOffset(msg.Offset, "")
		case err := <-checkingConsumer.Errors():
			log.Error("Failed to pull account-checking message from Kafka", "err", err)
			panic(err)
		}
	}
}
func T2est_sarama_msg_notwork(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	urls := []string{"192.168.112.32:9092"}
	topic := "platon-account-checking"

	consumer, err := sarama.NewConsumer(urls, consumerSaramaConfig())
	if err != nil {
		log.Error("Failed to init msg Kafka account-checking consumer....", "err", err)
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}()

	checkingConsumer, err := consumer.ConsumePartition(topic, int32(0), sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := checkingConsumer.Close(); err != nil {
			log.Error("Failed to close checkingConsumer", "err", err)
		}
	}()

	log.Info("Success to init msg Kafka account-checking consumer....")
	for {
		select {
		case msg := <-checkingConsumer.Messages():
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message", "offset", msg.Offset, "key", key, "value", value)

			//手工提交offset()
			/*if err := checkingConsumer.CommitOffsets(); err != nil {
				log.Error("Failed to commit checking consumer offset", "err", err)
				panic(err)
			}*/
		case err := <-checkingConsumer.Errors():
			log.Error("Failed to pull account-checking message from Kafka", "err", err)
			panic(err)
		}
	}
}

func clusterConfig() *cluster.Config {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_5_0_0
	//手工提交offset
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	//初始从最新的offset开始
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	return config
}

func T2est_sarama_cluster_msg_notwork(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	urls := []string{"192.168.112.32:9092"}
	topic := "platon-account-checking"
	if len(topic) == 0 {
		topic = defaultKafkaAccountCheckingTopic
	}

	consumer, err := cluster.NewConsumer(urls, defaultKafkaAccountCheckingConsumerGroup, []string{topic}, clusterConfig())
	if err != nil {
		log.Error("Failed to init msg Kafka account-checking consumer....", "err", err)
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}()

	log.Info("Success to init msg Kafka account-checking consumer....")
	for {
		select {
		case msg := <-consumer.Messages():
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message", "offset", msg.Offset, "key", key, "value", value)

			//手工提交offset()
			if err := consumer.CommitOffsets(); err != nil {
				log.Error("Failed to commit checking consumer offset", "err", err)
				panic(err)
			}
		case err := <-consumer.Errors():
			log.Error("Failed to pull account-checking message from Kafka", "err", err)
			panic(err)
		}
	}
}
