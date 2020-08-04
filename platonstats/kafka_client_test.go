package platonstats

import (
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

func Test_kafkaClient_consumer(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	kafkaClient := NewConfluentKafkaClient("192.168.9.201:9092", "", "platon-account-checking", "platon-account-checking-group")

	for {
		msg, err := kafkaClient.consumer.ReadMessage(-1)
		if err == nil {
			key := string(msg.Key)
			value := string(msg.Value)
			log.Debug("received account-checking message by group consumer", "key", key, "value", value)
		} else {
			// The client will automatically try to recover from all errors.
			log.Error("Consumer error", "msg", msg, "err", err)
		}
	}

}
