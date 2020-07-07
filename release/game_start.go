package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: game.start handler
func GameStart(c *gin.Context) {
	util.Info("GameStart ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get openId & roomId From Param
	var req model.GameStartReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "GameStart", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId

	// Find Room In Room Cache
	//_, ok := global.RoomCache[rid]
	//if !ok {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	ok, err := gocache.CheckRoomExists(rid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if !ok {
		util.Error("用户不在房间:%v", rid)
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find Room in RoomBucket
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameStart", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameStart Decoding Room Info ERROR[%v]", de.Error())
	//}
	err = gocache.GetRoomInfo(rid, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You are not in this room
	exist := false
	for _, v := range roomInfo.UnionIdSlice {
		if v == uid {
			exist = true
		}
	}
	if !exist {
		util.Error("用户不在房间:%v", rid)
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	// Find Offline User & Help Offline User Set Role
	existOffline := false
	users, err := gocache.GetRoomUsers(rid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_BELONG, c)
		return
	}
	//for _, p := range global.RoomCache[rid] {
	for _, p := range users {
		//for _, p := range uss.([]string) {
		online := false
		for _, o := range roomInfo.UnionIdSlice {
			if p == o {
				online = true
				break
			}
		}
		if !online {
			util.Debug("unionid:%v", p)
			for k, v := range roomInfo.PlayerSlice {
				if v.UnionId == "" {
					roomInfo.PlayerSlice[k].UnionId = p
					break
				}
			}
			existOffline = true
		}
	}

	// If Everyone has a Role
	for _, v := range roomInfo.PlayerSlice {
		if v.UnionId == "" && v.Role.Choice {
			util.Error("ERR_SOMEONE_NOT_SELECT")
			ResponseErr(model.ERR_SOMEONE_NOT_SELECT, c)
			return
		}
	}
	util.Debug("existOffline:%v", existOffline)
	// Update Room Bucket
	if existOffline {
		//encodingRoomInfo, err := json.Marshal(roomInfo)
		//if err != nil {
		//	//util.Logger(util.ERROR_LEVEL, "GameStart", "RoomInfo Encoding Err:"+err.Error())
		//	util.Error("GameStart RoomInfo Encoding ERROR[%v]", err.Error())
		//}
		//boltdb.CreateOrUpdate([]byte(rid), encodingRoomInfo, "RoomBucket")
		err = gocache.SetRoomInfo(rid, roomInfo)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ENTERED_ROOM, c)
			return
		}
	}

	// Send Broadcast Msg
	go websocket.SendGameStartMessage(roomInfo.UnionIdSlice)
	util.Info("success")

	//更新房间状态
	param := model.RoomRecordBase{}
	param.ScriptId = roomInfo.ScriptId
	param.RoomId = rid
	param.Owner = roomInfo.Owner
	param.Status = 1 //0创建房间1开始游戏2搜证结束3投票结束5答题结束6评分结束7游戏结束
	go point.RoomStatusUpdate(param)

	// Return Success
	resp := model.GameStartResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, resp)
}
