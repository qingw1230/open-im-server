package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/qingw1230/studyim/pkg/base_info"
	"github.com/qingw1230/studyim/pkg/common/config"
	"github.com/qingw1230/studyim/pkg/common/log"
	rpc "github.com/qingw1230/studyim/pkg/proto/auth"
	"github.com/qingw1230/studyim/pkg/proto/sdkws"
	"github.com/qingw1230/studyim/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserRegister(c *gin.Context) {
	params := api.UserRegisterReq{}
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

	req := &rpc.UserRegisterReq{UserInfo: &sdkws.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	req.OperationID = params.OperationID
	log.Info(req.OperationID, "UserRegister args ", req.String())
	etcdConn, err := grpc.NewClient("127.0.0.1:10600", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserRegister(context.Background(), req)
	if err != nil || reply.CommonResp.ErrCode != 0 {
		log.Error(req.OperationID, "UserRegister failed ", err, reply.CommonResp.ErrCode)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": reply.CommonResp.ErrMsg})
		return
	}

	pbDataToken := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	replyToken, err := client.UserToken(context.Background(), pbDataToken)
	if err != nil {
		log.Error(req.OperationID, "UserToken failed ", err.Error(), pbDataToken)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.UserRegisterResp{
		CommResp:  api.CommResp{ErrCode: replyToken.CommonResp.ErrCode, ErrMsg: replyToken.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.UserInfo.UserID, Token: replyToken.Token, ExpiredTime: replyToken.ExpiredTime},
	}
	log.Info(req.OperationID, "UserRegister return ", resp)
	c.JSON(http.StatusOK, resp)
}

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
	log.Info(req.OperationID, "UserToken args ", req.String())
	etcdConn, err := grpc.NewClient("127.0.0.1:10600", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
	}
	client := rpc.NewAuthClient(etcdConn)
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
