package auth

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/qingw1230/studyim/pkg/common/config"
	"github.com/qingw1230/studyim/pkg/common/constant"
	"github.com/qingw1230/studyim/pkg/common/db/mysql_model/im_mysql_model"
	"github.com/qingw1230/studyim/pkg/common/log"
	"github.com/qingw1230/studyim/pkg/common/token_verify"
	"github.com/qingw1230/studyim/pkg/grpc-etcdv3/getcdv3"
	pbAuth "github.com/qingw1230/studyim/pkg/proto/auth"
	"github.com/qingw1230/studyim/pkg/utils"
	"google.golang.org/grpc"
)

func (rpc *rpcAuth) UserToken(_ context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	_, err := im_mysql_model.GetUserByUserID(req.FromUserID)
	if err != nil {
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	token, expTime, err := token_verify.CreateToken(req.FromUserID, req.Platform)
	if err != nil {
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: token, ExpiredTime: expTime}, nil
}

type rpcAuth struct {
	pbAuth.UnimplementedAuthServer
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewRpcAuthServer(port int) *rpcAuth {
	log.NewPrivateLog("auth")
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcAuth) Run() {
	log.Info("0", "rpc auth start...")

	address := utils.ServerIP + ":" + strconv.Itoa(rpc.rpcPort)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("0", "listen network failed ", err.Error(), address)
		return
	}

	server := grpc.NewServer()
	defer server.GracefulStop()

	pbAuth.RegisterAuthServer(server, rpc)
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), utils.ServerIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		return
	}
	err = server.Serve(ln)
	if err != nil {
		return
	}
	log.Info("0", "rpc auth ok")
}
