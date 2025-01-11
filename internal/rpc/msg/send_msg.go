package msg

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/qingw1230/studyim/pkg/common/constant"
	"github.com/qingw1230/studyim/pkg/common/log"
	pbChat "github.com/qingw1230/studyim/pkg/proto/chat"
	"github.com/qingw1230/studyim/pkg/proto/sdkws"
	"github.com/qingw1230/studyim/pkg/utils"
)

func (rpc *rpcChat) encapsulateMsgData(msg *sdkws.MsgData) {
	msg.ServerMsgID = GetMsgID(msg.SendID)
	msg.SendTime = utils.GetCurrentTimestampByMill()
}

func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	log.Debug(pb.OperationID, "rpc SendMsg come here", pb.String())
	reply := pbChat.SendMsgResp{}
	rpc.encapsulateMsgData(pb.MsgData)
	msgToMQ := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID}

	switch pb.MsgData.SessionType {
	case constant.SingleChatType:
		msgToMQ.MsgData = pb.MsgData
		err1 := rpc.sendMsgToKafka(&msgToMQ, msgToMQ.MsgData.RecvID)
		if err1 != nil {
			log.Error(msgToMQ.OperationID, "kafka send msg err:RecvID", msgToMQ.MsgData.RecvID, msgToMQ.String())
			return returnMsg(&reply, pb, 201, "kafka send msg err", "", 0)
		}

		if msgToMQ.MsgData.SendID != msgToMQ.MsgData.RecvID {
			err2 := rpc.sendMsgToKafka(&msgToMQ, msgToMQ.MsgData.SendID)
			if err2 != nil {
				log.Error(msgToMQ.OperationID, "kafka send msg err:SendID", msgToMQ.MsgData.SendID, msgToMQ.String())
				return returnMsg(&reply, pb, 201, "kafka send msg err", "", 0)
			}
		}
		return returnMsg(&reply, pb, 0, "", msgToMQ.MsgData.ServerMsgID, msgToMQ.MsgData.SendTime)
	default:
		return returnMsg(&reply, pb, 203, "unkonwn sessionType", "", 0)
	}
}

func (rpc *rpcChat) sendMsgToKafka(m *pbChat.MsgDataToMQ, key string) error {
	pid, offset, err := rpc.producer.SendMessage(m, key)
	if err != nil {
		log.Error("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key)
	}
	return err
}

func GetMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t + "-" + sendID + "-" + strconv.Itoa(rand.Int())
}

func returnMsg(replay *pbChat.SendMsgResp, pb *pbChat.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*pbChat.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.SendTime = sendTime
	return replay, nil
}
