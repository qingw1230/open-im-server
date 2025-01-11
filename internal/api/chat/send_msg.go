package chat

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qingw1230/studyim/pkg/common/log"
	pbChat "github.com/qingw1230/studyim/pkg/proto/chat"
	im_sdk "github.com/qingw1230/studyim/pkg/proto/sdkws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paramsUserSendMsg struct {
	SenderPlatformID int32  `json:"senderPlatformID" binding:"required"`
	SendID           string `json:"sendID" binding:"required"`
	SenderNickName   string `json:"senderNickName"`
	SenderFaceURL    string `json:"senderFaceUrl"`
	OperationID      string `json:"operationID" binding:"required"`
	Data             struct {
		SessionType int32           `json:"sessionType" binding:"required"`
		MsgFrom     int32           `json:"msgFrom" binding:"required"`
		ContentType int32           `json:"contentType" binding:"required"`
		RecvID      string          `json:"recvID" `
		GroupID     string          `json:"groupID" `
		ForceList   []string        `json:"forceList"`
		Content     []byte          `json:"content" binding:"required"`
		Options     map[string]bool `json:"options" `
		ClientMsgID string          `json:"clientMsgID" binding:"required"`
		CreateTime  int64           `json:"createTime" binding:"required"`

		OffLineInfo *im_sdk.OfflinePushInfo `json:"offlineInfo" `
	}
}

func newUserSendMsgReq(token string, params *paramsUserSendMsg) *pbChat.SendMsgReq {
	pbData := pbChat.SendMsgReq{
		Token:       token,
		OperationID: params.OperationID,
		MsgData: &im_sdk.MsgData{
			SendID:           params.SendID,
			RecvID:           params.Data.RecvID,
			GroupID:          params.Data.GroupID,
			ClientMsgID:      params.Data.ClientMsgID,
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickName,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.Data.SessionType,
			MsgFrom:          params.Data.MsgFrom,
			ContentType:      params.Data.ContentType,
			Content:          params.Data.Content,
			CreateTime:       params.Data.CreateTime,
			Options:          params.Data.Options,
			OfflinePushInfo:  params.Data.OffLineInfo,
		},
	}
	return &pbData
}

func SendMsg(c *gin.Context) {
	params := paramsUserSendMsg{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		log.Error("json unmarshal err", "err", err.Error())
		return
	}

	token := c.Request.Header.Get("token")

	pbData := newUserSendMsgReq(token, &params)

	etcdConn, err := grpc.NewClient("127.0.0.1:10300", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
	}
	client := pbChat.NewChatClient(etcdConn)

	reply, err := client.SendMsg(context.Background(), pbData)
	if err != nil {
		log.Error(params.OperationID, "SendMsg rpc failed", params, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "SendMsg rpc failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
		"data": gin.H{
			"clientMsgID": reply.ClientMsgID,
			"serverMsgID": reply.ServerMsgID,
			"sendTime":    reply.SendTime,
		},
	})
}
