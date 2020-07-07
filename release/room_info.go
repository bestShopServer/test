package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: room.info handler
func RoomInfo(c *gin.Context) {
	util.Info("RoomInfo ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get ScriptId & OpenId From Param
	var req model.RoomInfoReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "RoomInfo", err, optionalFields) {
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
		util.Error("ERR_ROOM_NOT_EXIST")
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Find KV in Room Bucket
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	rif := model.RoomInfo{}
	//de := json.Unmarshal(b, &rif)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "WxJoin", "Decoding Room Info Err:"+de.Error())
	//	util.Error("WxJoin Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(rid, &rif)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You are not in this room
	exist := false
	for _, v := range rif.UnionIdSlice {
		if v == uid {
			exist = true
		}
	}
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	// Return roomId & playerSlice
	resp := model.RoomInfoResp{}
	resp.Params.RoomId = rid
	resp.Params.RoomInfo = rif
	util.Info("success:%+v", resp)
	c.JSON(http.StatusOK, resp)
}
