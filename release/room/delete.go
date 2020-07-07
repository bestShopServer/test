package room

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
)

//房间删除
func RoomDelete(c *gin.Context) {
	util.Info("RoomDelete ...")

	var req model.RoomDeleteReq

	err := c.Bind(&req)
	util.Info("请求参数[%+v]", req)

	// Optional Fields List
	optionalFields := []string{}
	// Check Param
	if !release.CheckParams(req, "RoomDelete", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	//判断房间是否存在
	ok, err := gocache.CheckRoomExists(req.RoomId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if !ok {
		util.Error("ERR_ROOM_NOT_EXIST")
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	roomInfo := model.RoomInfo{}
	err = gocache.GetRoomInfo(req.RoomId, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	//判断房主
	if req.UnionId != roomInfo.Owner {
		util.Error("用户:%v不是房主:%v", req.UnionId, req.RoomId)
		release.ResponseErr(model.ERR_NOT_ROOM_OWNER, c)
		return
	}

	//删除记录用户房间信息
	rooms, err := gocache.GetRoomUsers(req.RoomId)
	for _, uid := range rooms {
		roomCode, err := gocache.GetUserRoom(uid)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			release.ResponseErr(model.ERR_ROOM_DELETE, c)
			return
		}
		//用户在当前房间的删除记录
		if roomCode == req.RoomId {
			ok, err = gocache.Delete(uid)
			if err != nil {
				util.Error("ERROR[%v]", err.Error())
				release.ResponseErr(model.ERR_ROOM_DELETE, c)
				return
			}
			if !ok {
				util.Error("删除房间失败[%v]", ok)
				release.ResponseErr(model.ERR_ROOM_DELETE, c)
				return
			}
		}
	}

	//删除房间信息
	ok, err = gocache.Delete(req.RoomId)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		release.ResponseErr(model.ERR_ROOM_DELETE, c)
		return
	}
	if !ok {
		util.Error("删除房间失败[%v]", ok)
		release.ResponseErr(model.ERR_ROOM_DELETE, c)
		return
	}

	//长链通知用户房间已删除
	go websocket.SendRoomDeleteMessage(roomInfo.UnionIdSlice)

	Resp := model.ErrResp{}
	release.ResponseOk(c, &Resp)
}
