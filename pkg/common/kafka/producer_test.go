package kafka

import (
	"testing"

	"github.com/qingw1230/studyim/pkg/proto/auth"
)

func Test_SendMessage(t *testing.T) {
	p := NewKafkaProducer([]string{"127.0.0.1:9092"}, "test")
	msg := auth.UserTokenReq{
		Platform:   5,
		FromUserID: "test-qgw",
	}
	p.SendMessage(&msg, "testID")
}
