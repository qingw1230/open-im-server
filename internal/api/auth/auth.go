package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/qingw1230/studyim/pkg/base_info"
	"github.com/qingw1230/studyim/pkg/common/config"
	"github.com/qingw1230/studyim/pkg/common/log"
	"github.com/qingw1230/studyim/pkg/discoveryregistry/zookeeper"
	rpc "github.com/qingw1230/studyim/pkg/proto/auth"
)

func UserToken(c *gin.Context) {
	params := api.UserTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		log.Error("0", "BindJSON failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if params.Secret != config.Config.Secret {
		log.Error(params.OperationID, "params.Secret != config.Config.Secret", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "not authorized"})
		return
	}

	req := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	conn, _ := zookeeper.ZK.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImAuthName)
	client := rpc.NewAuthClient(conn)
	reply, err := client.UserToken(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp := api.UserTokenResp{
		CommResp:  api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.FromUserID, Token: reply.Token, ExpiredTime: reply.ExpiredTime},
	}
	c.JSON(http.StatusOK, resp)
}
