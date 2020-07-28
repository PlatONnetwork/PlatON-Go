package platonstats

import (
	"strings"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/Shopify/sarama"
)

const (
	// historyUpdateRange is the number of blocks a node should report upon login or
	// history request.
	sampleEventChanSize                      = 50
	defaultKafkaBlockTopic                   = "platon-block"
	defaultKafkaAccountCheckingConsumerGroup = "platon-account-checking-group"
	defaultKafkaAccountCheckingTopic         = "platon-account-checking"
)

type KafkaClient struct {
	brokers                []string
	blockTopic             string
	accountCheckingTopic   string
	saramaClient           sarama.Client
	asyncProducer          sarama.AsyncProducer
	syncProducer           sarama.SyncProducer
	consumer               sarama.Consumer
	consumerGroup          sarama.ConsumerGroup
	partitionConsumer      sarama.PartitionConsumer
	offsetManager          sarama.OffsetManager
	partitionOffsetManager sarama.PartitionOffsetManager
}

func producerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.MaxMessageBytes = 500000000
	return config
}

func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Return.Errors = true

	//手工提交offset
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	//初始从最新的offset开始
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}

func saramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0

	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.MaxMessageBytes = 500000000

	config.Consumer.Return.Errors = true

	//手工提交offset
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	//初始从最新的offset开始
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}

func NewKafkaClient(urls, blockTopic, checkingTopic string) *KafkaClient {
	brokers := strings.Split(urls, ",")

	if len(blockTopic) == 0 {
		blockTopic = defaultKafkaAccountCheckingTopic
	}

	if len(checkingTopic) == 0 {
		checkingTopic = defaultKafkaAccountCheckingTopic
	}

	kafkaClient := &KafkaClient{
		brokers:              brokers,
		blockTopic:           blockTopic,
		accountCheckingTopic: checkingTopic,
	}

	if blockProducer, err := sarama.NewSyncProducer(brokers, producerConfig()); err != nil {
		log.Error("Failed to create Kafka producer")
		panic(err)
	} else {
		kafkaClient.syncProducer = blockProducer
	}

	if consumerGroup, err := sarama.NewConsumerGroup(brokers, defaultKafkaAccountCheckingConsumerGroup, consumerConfig()); err != nil {
		log.Error("Failed to create Kafka consumer")
		panic(err)
	} else {
		kafkaClient.consumerGroup = consumerGroup
		/*if partitionConsumer, err := consumer.ConsumePartition(checkingTopic, 0, sarama.OffsetNewest); err != nil {
			log.Error("Failed to create Kafka partition consumer")
			panic(err)
		} else {
			kafkaClient.partitionConsumer = partitionConsumer
		}*/
	}
	return kafkaClient
}

func NewKafkaClie3nt(urls, blockTopic, checkingTopic string) *KafkaClient {
	brokers := strings.Split(urls, ",")

	if len(blockTopic) == 0 {
		blockTopic = defaultKafkaAccountCheckingTopic
	}

	if len(checkingTopic) == 0 {
		checkingTopic = defaultKafkaAccountCheckingTopic
	}

	kafkaClient := &KafkaClient{
		brokers:              brokers,
		blockTopic:           blockTopic,
		accountCheckingTopic: checkingTopic,
	}

	if blockProducer, err := sarama.NewSyncProducer(brokers, producerConfig()); err != nil {
		log.Error("Failed to create Kafka producer")
		panic(err)
	} else {
		kafkaClient.syncProducer = blockProducer
	}

	if consumer, err := sarama.NewConsumer(brokers, nil); err != nil {
		log.Error("Failed to create Kafka consumer")
		panic(err)
	} else {
		kafkaClient.consumer = consumer
		if partitionConsumer, err := consumer.ConsumePartition(checkingTopic, 0, sarama.OffsetNewest); err != nil {
			log.Error("Failed to create Kafka partition consumer")
			panic(err)
		} else {
			kafkaClient.partitionConsumer = partitionConsumer
		}
	}
	return kafkaClient
}
func NewKafkaClien2t(urls, blockTopic, checkingTopic string) *KafkaClient {
	brokers := strings.Split(urls, ",")

	if len(blockTopic) == 0 {
		blockTopic = defaultKafkaAccountCheckingTopic
	}

	if len(checkingTopic) == 0 {
		checkingTopic = defaultKafkaAccountCheckingTopic
	}

	client, err := sarama.NewClient(brokers, saramaConfig())
	if err != nil {
		log.Error("Failed to create Kafka client")
		panic(err)
	}

	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Error("kafka connect error:")
		panic(err)
	}

	syncProducer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Error("kafka connect error:")
		panic(err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Error("kafka connect error:")
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(checkingTopic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Error("Failed to create Kafka partition_consumer....", "err", err)
		panic(err)
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient(defaultKafkaAccountCheckingConsumerGroup, client)
	if err != nil {
		panic("offsetManager create error")
	}

	partitionOffsetManager, err := offsetManager.ManagePartition(checkingTopic, int32(0))
	if err != nil {
		panic("offsetManager create error")
	}

	return &KafkaClient{
		brokers:                brokers,
		blockTopic:             blockTopic,
		accountCheckingTopic:   checkingTopic,
		asyncProducer:          asyncProducer,
		syncProducer:           syncProducer,
		consumer:               consumer,
		partitionConsumer:      partitionConsumer,
		offsetManager:          offsetManager,
		partitionOffsetManager: partitionOffsetManager,
	}
}

func (kc *KafkaClient) Close() {
	if kc.saramaClient != nil {
		if err := kc.saramaClient.Close(); err != nil {
			log.Error("Failed to close Kafka client", "err", err)
		}
	}
	if kc.asyncProducer != nil {
		if err := kc.asyncProducer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}
	if kc.syncProducer != nil {
		if err := kc.syncProducer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}
	if kc.consumer != nil {
		if err := kc.consumer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}
	if kc.partitionConsumer != nil {
		if err := kc.partitionConsumer.Close(); err != nil {
			log.Error("Failed to close partition consumer", "err", err)
		}
	}
	if kc.offsetManager != nil {
		if err := kc.offsetManager.Close(); err != nil {
			log.Error("Failed to close offset manager", "err", err)
		}
	}
	if kc.partitionOffsetManager != nil {
		if err := kc.partitionOffsetManager.Close(); err != nil {
			log.Error("Failed to close partition offset manager", "err", err)
		}
	}
}
