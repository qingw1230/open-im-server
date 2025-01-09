package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/qingw1230/studyim/pkg/common/log"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	addr     []string
	topic    string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{
		addr:  addr,
		topic: topic,
	}
	p.config = sarama.NewConfig()
	p.config.Producer.Return.Successes = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll
	p.config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(p.addr, p.config)
	if err != nil {
		panic(err.Error())
	}
	p.producer = producer
	return &p
}

// SendMessage 向 kafka 发送消息
// SingleChatType key 为接收者的 userID
func (p *Producer) SendMessage(m proto.Message, key ...string) (int32, int64, error) {
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	if len(key) == 1 {
		kMsg.Key = sarama.StringEncoder(key[0])
	}
	bMsg, err := proto.Marshal(m)
	if err != nil {
		log.Error("", "proto marshal err: %s", err.Error())
		return -1, -1, err
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	return p.producer.SendMessage(kMsg)
}
