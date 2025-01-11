package msg

import (
	"net"
	"strconv"
	"strings"

	"github.com/qingw1230/studyim/pkg/common/config"
	"github.com/qingw1230/studyim/pkg/common/kafka"
	"github.com/qingw1230/studyim/pkg/common/log"
	"github.com/qingw1230/studyim/pkg/grpc-etcdv3/getcdv3"
	pbChat "github.com/qingw1230/studyim/pkg/proto/chat"
	"github.com/qingw1230/studyim/pkg/utils"
	"google.golang.org/grpc"
)

type rpcChat struct {
	pbChat.UnimplementedChatServer
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	producer        *kafka.Producer
}

func NewRpcChatServer(port int) *rpcChat {
	log.NewPrivateLog("msg")
	rpc := rpcChat{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOfflineMessageName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	rpc.producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	return &rpc
}

func (rpc *rpcChat) Run() {
	address := utils.ServerIP + ":" + strconv.Itoa(rpc.rpcPort)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("", "", "listen network failed, err: %s, address: %s", err.Error(), address)
		return
	}

	server := grpc.NewServer()
	defer server.GracefulStop()

	pbChat.RegisterChatServer(server, rpc)
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), utils.ServerIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "", "register rpc failed, err: %s", err.Error())
		return
	}

	err = server.Serve(ln)
	if err != nil {
		return
	}
}
