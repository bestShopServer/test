package user

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
)

//投降
func UserSurrender(c *gin.Context) {
	util.Info("WxUserRoomBase ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Get Params
	var req model.UserSurrenderReq
	err := c.BindJSON(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Params
	if !release.CheckParams(req, "UserSurrender", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	//判断用户是否在此房间
	cache_room_id, err := gocache.GetUserRoom(req.UnionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_GAME_HAS_OVER, c)
		return
	}
	util.Debug("cache_room_id:%v", cache_room_id)

	// If user cache not eq room id
	if cache_room_id != req.RoomId {
		util.Error("ERR_ROOM_LINK")
		release.ResponseErr(model.ERR_ROOM_LINK, c)
		return
	}

	//查看房间状态是否已结束
	roomInfo := model.RoomInfo{}
	err = gocache.GetRoomInfo(req.RoomId, &roomInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}
	//记录投降用户
	for i, tmp := range roomInfo.UserInfoSlice {
		if tmp.UnionId == req.UnionId {
			tmp.Surrender = true
			release.UserInfoSliceModify(&roomInfo.UserInfoSlice, i, tmp)
			break
		}
	}

	//在线人数-1就可以结束本局游戏
	num := len(roomInfo.UnionIdSlice)
	util.Info("房间剩余人数:%v", num)
	for _, uid := range roomInfo.UnionIdSlice {
		for _, tmp := range roomInfo.UserInfoSlice {
			if tmp.UnionId == uid && tmp.Surrender {
				num -= 1
			}
		}
	}
	util.Info("房间未投降人数:%v", num)

	if num <= 1 { //剩余1人未投降，则房间游戏结束
		//房间状态更新为游戏结束
		roomInfo.Status = 2 //游戏结束

		for _, tmp := range roomInfo.UserInfoSlice {
			//记录游戏得分为0
			err = gocache.SetRoomUserQuestionScore(req.RoomId, tmp.UnionId, 0, 0)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				release.ResponseErr(model.ERR_GET_SCRIPTS, c)
				return
			}
		}

		//长链通知用户房间已结束，可以直接看解析
		util.Info("通知用户:%+v 房间游戏结束!", roomInfo.UnionIdSlice)
		go websocket.SendRoomGameEndMessage(roomInfo.UnionIdSlice)
	}

	//记录房间数据
	err = gocache.SetRoomInfo(req.RoomId, roomInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}

	// Return
	release.ResponseSuccess(c)
}
