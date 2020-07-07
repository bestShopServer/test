package release

//
//import (
//	"DetectiveMasterServer/boltdb"
//	"DetectiveMasterServer/global"
//	"DetectiveMasterServer/model"
//	"DetectiveMasterServer/util"
//	"DetectiveMasterServer/websocket"
//	"github.com/gin-gonic/gin"
//	"github.com/json-iterator/go"
//	"net/http"
//)
//
//// Func: room.enter handler
//func RoomEnter(c *gin.Context) {
//	util.Info("RoomEnter ...")
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	var user model.UserInfo
//
//	// Get openId & roomId From Param
//	var req model.RoomEnterReq
//	err := c.Bind(&req)
//
//	// Optional Fields List
//	var optionalFields []string
//
//	// Check Param
//	if !CheckParams(req, "RoomEnter", err, optionalFields) {
//		util.Error("ERR_WRONG_FORMAT")
//		ResponseErr(model.ERR_WRONG_FORMAT, c)
//		return
//	}
//
//	//uid := req.OpenId
//	uid := req.UnionId
//	rid := req.RoomId
//	user.UnionId = req.UnionId
//
//	// Find Room In Room Cache
//	_, ok := global.RoomCache[rid]
//	if !ok {
//		util.Error("ERR_ROOM_NOT_EXIST")
//		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
//		return
//	}
//
//	// Find KV in Room Bucket
//	b := boltdb.View([]byte(rid), "RoomBucket")
//	if b == nil {
//		util.Error("ERR_ROOM_NOT_EXIST")
//		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
//		return
//	}
//
//	// Decode Room Info
//	rif := model.RoomInfo{}
//	de := json.Unmarshal(b, &rif)
//	if de != nil {
//		//util.Logger(util.ERROR_LEVEL, "RoomEnter", "Decoding Room Info Err:"+de.Error())
//		util.Error("RoomEnter Decoding Room Info ERROR[%v]", de.Error())
//	}
//
//	// If openIdSlice length eq playerSlice
//	if len(rif.UnionIdSlice) == len(rif.PlayerSlice) {
//		util.Error("ERR_ROOM_IS_FULL")
//		ResponseErr(model.ERR_ROOM_IS_FULL, c)
//		return
//	}
//
//	// If openIdSlice Contains openId
//	for _, v := range rif.UnionIdSlice {
//		if v == uid {
//			util.Error("ERR_ENTERED_ROOM")
//			ResponseErr(model.ERR_ENTERED_ROOM, c)
//			return
//		}
//	}
//
//	// Add openId to openIdSlice
//	rif.UnionIdSlice = append(rif.UnionIdSlice, uid)
//
//	// Update BoltDB
//	roomInfoEncoding, err := json.Marshal(rif)
//	if err != nil {
//		//util.Logger(util.ERROR_LEVEL, "RoomEnter", "Room Info Encoding Err:"+err.Error())
//		util.Error("RoomEnter Room Info Encoding ERROR[%v]", err.Error())
//	}
//	boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
//
//	// Update UserCache & RoomCache
//	global.SetUserCache(uid, rid)
//	global.AddUserToRoomCache(uid, rid)
//
//	// Send Room Enter Msg
//	//go websocket.SendRoomEnterMessage(rif.UnionIdSlice, uid)
//	go websocket.SendRoomEnterMessage(rif.UnionIdSlice, user)
//	util.Info("success")
//
//	// Return Room Info
//	resp := model.RoomEnterResp{}
//	resp.Params.RoomId = rid
//	resp.Params.RoomInfo = rif
//	c.JSON(http.StatusOK, resp)
//
//}
