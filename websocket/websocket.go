package websocket

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
	"net/url"
	"sync"
)

type WebSocketConn struct {
	Conn     *websocket.Conn
	SendChan chan interface{}
	//RecvChan chan []byte
	Mutex sync.Mutex
}

type RedisChannel struct {
	Conn     redis.Conn
	SendChan chan interface{}
	Mutex    sync.Mutex
}

//var WS *websocket.Conn
var WsConn WebSocketConn
var conf = config.GetConfig()
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//var RcSendConn RedisChannel
var RcRecvConn RedisChannel

// WebSocket Request Enum Type
const (
	HAND_SHAKE_REQ = iota
	USER_LOGIN_REQ
	BROAD_CAST_REQ
)

const (
	USER_OFFLINE_RESP = iota
)

// WebSocket Log Category
const WS_CATEGORY = "WebSocket"

// Common Message
type WebSocketMsg struct {
	ReqType int `json:"req_type"`
	//OpenId  string `json:"open_id"`
	UnionId string `json:"union_id"`
}

// Handshake Message
type HandShakeMsg struct {
	WebSocketMsg
}

// User Login Message
type UserLoginMsg struct {
	WebSocketMsg
	Message string `json:"message"`
}

// Broadcast Message
type BroadCastMsg struct {
	WebSocketMsg
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type WebSocketResp struct {
	RespType int    `json:"resp_type"`
	Message  string `json:"message"`
}

// Func: Connect WebSocket Server
//func DialWebSocket() *websocket.Conn {
func DialWebSocket() WebSocketConn {

	u := url.URL{
		Scheme: "ws",
		Host:   conf.WebSocketHost,
		Path:   conf.WebSocketPath,
	}
	//util.Logger(util.INFO_LEVEL, WS_CATEGORY, "Connecting to "+u.String())
	util.Info("Connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		//util.Logger(util.ERROR_LEVEL, WS_CATEGORY, err.Error())
		util.Error("WebSocket ERROR[%v]", err.Error())
	}
	//defer c.Close()

	if c != nil {
		//util.Debug("conn c:%+v", c)
		WsConn = WebSocketConn{
			Conn:     c,
			SendChan: make(chan interface{}, 128),
		}
		//go ReadMessage()
		go WsConn.ReadMessage()
	}
	//util.Debug("conn c:%+v", c)
	util.Info("Connect WebSocket OK! %+v", WsConn)

	return WsConn
}

// Func: Reconnect To WebSocket Server
func ReConnectWebSocket() bool {
	//WS = DialWebSocket()
	//if WS == nil {
	WsConn = DialWebSocket()
	if WsConn.Conn == nil {
		return false
	} else {
		SendMessage(HAND_SHAKE_REQ, conf.ServerOpenId, "", "")
		return true
	}
}

// Func: Write Message
//func WriteMessage(text interface{}) {
func (conn *WebSocketConn) WriteMessage() {
	conn.Mutex.Lock()
	//for {
	msg := <-conn.SendChan
	err := conn.Conn.WriteJSON(msg)
	if err != nil {
		//util.Logger(util.ERROR_LEVEL, WS_CATEGORY, "Sending Msg Err: "+err.Error())
		util.Error("WebSocket Sending Msg ERROR[%v]", err.Error())
		return
	}

	logText, _ := json.Marshal(msg)
	//util.Logger(util.INFO_LEVEL, WS_CATEGORY, "Sending Msg: "+string(logText))
	util.Info("WebSocket Sending Msg: %v", string(logText))
	//}
	conn.Mutex.Unlock()
}

// Func: Read Message
func (conn *WebSocketConn) ReadMessage() {
	for {
		//_, message, err := WS.ReadMessage()
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			//util.Logger(util.ERROR_LEVEL, WS_CATEGORY, err.Error())
			util.Error("ReadMessage ERROR[%v]", err.Error())
			defer conn.Conn.Close()
			conn.Conn = nil
			return
		}

		//util.Logger(util.INFO_LEVEL, WS_CATEGORY, "Receiving Msg: "+string(message))
		util.Info("Receiving Msg:", string(message))
		go Transport(message)
	}
}

func Transport(t []byte) {
	ws := WebSocketResp{}
	err := json.Unmarshal(t, &ws)
	if err != nil {
		//util.Logger(util.ERROR_LEVEL, WS_CATEGORY, "Unmarshal WsMsg Err:"+err.Error())
		util.Error("Unmarshal WsMsg ERROR [%v]", err.Error())
		return
	}

	switch ws.RespType {
	case USER_OFFLINE_RESP:
		go DealUserOffline(ws.Message)
	default:
		return
	}
}

// Func: Send Message To WebSocket Server
func SendMessage(req int, unionId string, message string, detail string) {
	var text interface{}
	var msg []byte
	var err error
	wsm := WebSocketMsg{
		//OpenId:  openId,
		UnionId: unionId,
		ReqType: req,
	}
	util.Info("UnionId[%v] ReqType[%v]", wsm.UnionId, wsm.ReqType)

	switch req {
	case HAND_SHAKE_REQ:
		text = &HandShakeMsg{
			wsm,
		}
	case USER_LOGIN_REQ:
		text = &UserLoginMsg{
			wsm,
			message,
		}
	case BROAD_CAST_REQ:
		text = &BroadCastMsg{
			wsm,
			message,
			detail,
		}
	default:
		text = &WebSocketMsg{}
	}

	//if WS == nil {
	//if WsConn.Conn == nil {
	//	if ReConnectWebSocket() {
	//		WriteMessage(text)
	//	}
	//} else {
	//	WriteMessage(text)
	//}

	//test redis
	//if WsConn.Conn == nil {
	//	if !ReConnectWebSocket() {
	//		util.Error("ReConnectWebSocket ERROR! union:%v", unionId)
	//		return
	//	}
	//}

	//WsConn.SendChan <- text
	//go WsConn.WriteMessage()
	util.Debug("publish %+v", text)
	msg, err = json.Marshal(text)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	//RcSendConn.SendChan <- msg
	//go RcSendConn.PushMessage()

	go RedisPublishMessage(msg)
}

// Func: Send Handshake Message
func SendHandShakeMessage() {
	SendMessage(HAND_SHAKE_REQ, conf.ServerOpenId, "", "")
}

// Func: Send User Login Message
func SendUserLoginMessage(unionId string) {
	SendMessage(USER_LOGIN_REQ, conf.ServerOpenId, unionId, "")
}

// Func: Send Room Enter Message
func SendRoomEnterMessage(unionIdList []string, use model.UserInfo) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitRoomEnterMsg()
	rem.Data = use
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Send Role Choose Message
func SendRoleChooseMessage(unionIdList []string, uid string, rid int) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitRoleChooseMsg()
	rem.Data.OpenId = uid
	rem.Data.RoleId = rid
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Send Game Start Message
func SendGameStartMessage(unionIdList []string) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitGameStartMsg()
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Send Room Exit Message
func SendRoomExitMessage(unionIdList []string, uid string) {
	openIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitRoomExitMsg()
	rem.Data = uid
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(openIdListEncoding), string(msgDetail))
}

// Func: Send Clew Open Message
func SendClewOpenMessage(unionIdList []string, unionId string, key string, index []int) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitClewOpenMsg()
	//rem.Data.OpenId = openId
	rem.Data.UnionId = unionId
	rem.Data.Key = key
	rem.Data.Index = index
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Send Clew Open Message
func SendClewPubMessage(unionIdList []string, unionId string, key string, index []int) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitClewPubMsg()
	//rem.Data.OpenId = openId
	rem.Data.UnionId = unionId
	rem.Data.Key = key
	rem.Data.Index = index
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Stage Next Message
func SendStageNextMsg(unionIdList []string) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitStageNextMsg()
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Step Into Game Vote Message
func SendStepIntoGameVoteMsg(unionIdList []string) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitStepIntoVoteMsg()
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Not Vote People Left Message
func SendGameVoteMsg(unionIdList []string, notVote int) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitGameVoteMsg()
	rem.Data = notVote
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

// Func: Deal User Offline Resp
func DealUserOffline(unionId string) {

	// Find RoomId From UserCache
	//roomId, ok := global.UserCache[unionId]
	//if !ok || roomId == "" {
	//	return
	//}
	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	//roomId, err := gocache.GetUserRoom(unionId)
	roomId, err := gocache.ConnGetUserRoom(conn, unionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return
	}

	// Remove Your OpenId From Room Bucket
	//rb := boltdb.View([]byte(roomId), "RoomBucket")
	//if rb != nil {
	//	rif := model.RoomInfo{}
	//	err := json.Unmarshal(rb, &rif)
	//
	//	if err != nil {
	//		util.Logger(util.ERROR_LEVEL, "GlobalCore", err.Error())
	//		return
	//	}
	//
	//	for k, v := range rif.UnionIdSlice {
	//		if v == unionId {
	//			rif.UnionIdSlice = append(rif.UnionIdSlice[:k], rif.UnionIdSlice[k+1:]...)
	//			break
	//		}
	//	}
	//
	//	// Update RoomBucket
	//	encodingRif, _ := json.Marshal(rif)
	//	boltdb.CreateOrUpdate([]byte(roomId), encodingRif, "RoomBucket")
	//
	//	// Send Room Exit Msg
	//	go SendRoomExitMessage(rif.UnionIdSlice, unionId)
	//
	//	// If You Are In Vote Stage Then Remove Your OpenId And Send Game Vote Message
	//	vb := boltdb.View([]byte(roomId), "VoteBucket")
	//	if vb != nil {
	//		var votes map[string]bool
	//		json.Unmarshal(vb, &votes)
	//		if votes[unionId] == false {
	//			votes[unionId] = true
	//			notVoteNum := global.CalcNotVoteNum(votes)
	//			encodingVotes, _ := json.Marshal(votes)
	//			boltdb.CreateOrUpdate([]byte(roomId), encodingVotes, "VoteBucket")
	//			go SendGameVoteMsg(rif.UnionIdSlice, notVoteNum)
	//		}
	//	}
	//}

	rif := model.RoomInfo{}
	//err = gocache.GetRoomInfo(roomId, &rif)
	err = gocache.ConnGetRoomInfo(conn, roomId, &rif)

	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return
	}
	//util.Debug("多幕:%+v", rif.PlayerSlice)

	for k, v := range rif.UnionIdSlice {
		if v == unionId {
			rif.UnionIdSlice = append(rif.UnionIdSlice[:k], rif.UnionIdSlice[k+1:]...)
			break
		}
	}

	// Update RoomBucket
	//util.Debug("房间 %v 数据:%+v", roomId, rif)
	//util.Debug("多幕:%+v", rif.PlayerSlice)
	//err = gocache.SetRoomInfo(roomId, rif)
	err = gocache.ConnSetRoomInfo(conn, roomId, rif)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//release.ResponseErr(model.ERR_ENTERED_ROOM, c)
		return
	}

	// Send Room Exit Msg
	go SendRoomExitMessage(rif.UnionIdSlice, unionId)

	// If You Are In Vote Stage Then Remove Your OpenId And Send Game Vote Message
	//暂时先删除掉线后自动投票
	//votes, bl, err := gocache.GetVoteInfo(roomId)
	//if err != nil {
	//	util.Error("ERROR:%v", err.Error())
	//	return
	//}
	//if bl {
	//	//var votes map[string]bool
	//	//json.Unmarshal(vb, &votes)
	//	if votes[unionId] == false {
	//		votes[unionId] = true
	//		notVoteNum := global.CalcNotVoteNum(votes)
	//		//encodingVotes, _ := json.Marshal(votes)
	//		//boltdb.CreateOrUpdate([]byte(roomId), encodingVotes, "VoteBucket")
	//		err = gocache.SetVoteInfo(roomId, votes)
	//		if err != nil {
	//			util.Error("ERROR:%v", err.Error())
	//			//release.ResponseErr(model.ERR_NOT_ALL_VOTED, c)
	//			return
	//		}
	//		go SendGameVoteMsg(rif.UnionIdSlice, notVoteNum)
	//	}
	//}

}

// Func: Send Room Delete Message 长链通知房间解散
func SendRoomDeleteMessage(unionIdList []string) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitRoomDeleteMsg()
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

//长链通知房间游戏结束
func SendRoomGameEndMessage(unionIdList []string) {
	unionIdListEncoding, _ := json.Marshal(unionIdList)

	rem := model.InitGameEndMsg()
	msgDetail, _ := json.Marshal(rem)

	SendMessage(BROAD_CAST_REQ, conf.ServerOpenId, string(unionIdListEncoding), string(msgDetail))
}

//用于同步长链推送消息
func DialRedisChannel() {
	//conn := gocache.RedisConnPool.Get()

	// Setup Initialize the Redis instance
	//util.Info("redis addr:%v", config.GetConfig().RedisAddr)
	//send_conn, err := redis.Dial("tcp", config.GetConfig().RedisAddr)
	//if len(config.GetConfig().RedisAuth) > 0 {
	//	if _, err := send_conn.Do("AUTH", config.GetConfig().RedisAuth); err != nil {
	//		//defer RedisConn.Close()
	//		util.Error("ERROR:%v", err.Error())
	//		return
	//	}
	//	util.Info("redis auth passwd")
	//}
	//_, err = send_conn.Do("select", config.GetConfig().RedisDb)
	//if err != nil {
	//	util.Error("ERROR:%v", err.Error())
	//	return
	//}
	//util.Info("redis select %v", config.GetConfig().RedisDb)

	util.Info("redis addr:%v", config.GetConfig().RedisAddr)
	recv_conn, err := redis.Dial("tcp", config.GetConfig().RedisAddr)
	if len(config.GetConfig().RedisAuth) > 0 {
		if _, err := recv_conn.Do("AUTH", config.GetConfig().RedisAuth); err != nil {
			//defer RedisConn.Close()
			util.Error("ERROR:%v", err.Error())
			return
		}
		util.Info("redis auth passwd")
	}
	_, err = recv_conn.Do("select", config.GetConfig().RedisDb)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return
	}
	util.Info("redis select %v", config.GetConfig().RedisDb)

	//RcSendConn = RedisChannel{
	//	Conn:     send_conn,
	//	SendChan: make(chan interface{}, 128),
	//}
	RcRecvConn = RedisChannel{
		Conn:     recv_conn,
		SendChan: make(chan interface{}, 128),
	}
	//go ReadMessage()
	go RcRecvConn.SubScribeMessage()

	//util.Debug("conn c:%+v", c)
	util.Info("Connect  Redis Channel OK!")

	return
}

//接收消息
func (conn *RedisChannel) SubScribeMessage() {
	//timeout := time.NewTimer(time.Second * 10)

	psc := redis.PubSubConn{conn.Conn}
	svrs := config.GetConfig().RedisRecvChannels
	for _, tmp := range svrs {
		go func() {
			err := psc.Subscribe(tmp)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
			}
			for {
				switch v := psc.Receive().(type) {
				case redis.Message:
					util.Info("subscribe: %s: message: %s\n", v.Channel, v.Data)
					go Transport(v.Data)
				case redis.Subscription:
					util.Info("%+v", v)
				case error:
					util.Error("ERROR:%+v", v)
				}
			}
		}()
	}
}

//发送消息
func (conn *RedisChannel) PushMessage() {
	svrs := config.GetConfig().RedisSendChannels
	msg := <-conn.SendChan
	for _, tmp := range svrs {
		//util.Debug("publish %v %v", tmp, msg)
		_, err := conn.Conn.Do("Publish", tmp, msg)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
		}
	}
}

//发送消息
func RedisPublishMessage(msg []byte) {
	svrs := config.GetConfig().RedisSendChannels
	conn := gocache.RedisConnPool.Get()
	for _, tmp := range svrs {
		//util.Debug("publish %v %v", tmp, msg)
		_, err := conn.Do("Publish", tmp, msg)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
		}
	}
}
