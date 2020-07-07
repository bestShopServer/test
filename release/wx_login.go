package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: wx.login handler
func WxLogin(c *gin.Context) {
	util.Info("WxLogin ...")
	// Get Code Param
	var req model.WxLoginReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "WxLogin", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)
	//openId := req.OpenId

	// Save OpenId in Database if it Not Exist
	//ob := boltdb.View([]byte(req.UnionId), "UserBucket")
	//if ob == nil {
	//	ct := time.Now().Format("2006-01-02 15:04:05")
	//	boltdb.CreateOrUpdate([]byte(req.UnionId), []byte(ct), "UserBucket")
	//}
	conn := gocache.RedisConnPool.Get()
	defer conn.Close()

	// Check If User In A Room Not Over Yet
	//roomId := global.UserCache[req.UnionId]
	//roomId, err := gocache.GetUserRoom(req.UnionId) //存储用户对应的房间号
	roomId, err := gocache.ConnGetUserRoom(conn, req.UnionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}
	//err = gocache.SetRoomUser(roomId, req.UnionId)
	err = gocache.ConnSetRoomUser(conn, roomId, req.UnionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_LINK, c)
		return
	}
	/* add by skc at 2020-04-29 begin */
	roomInfo := model.RoomInfo{}
	err = gocache.ConnGetRoomInfo(conn, roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_LINK, c)
		return
	}
	exist := false
	for _, uid := range roomInfo.UnionIdSlice {
		if uid == req.UnionId {
			util.Info("用户:%v在房间:%v", uid, roomId)
			exist = true
			break
		}
	}
	if !exist { //用户离开房间后重新进入房间
		roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice, req.UnionId) //添加用户
		err = gocache.ConnSetRoomInfo(conn, roomId, roomInfo)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ROOM_LINK, c)
			return
		}
		util.Info("用户:%v重新回房间:%v", req.UnionId, roomId)
	}
	/* add by skc at 2020-04-29 end */

	// Notify WebSocket User Login
	go websocket.SendUserLoginMessage(req.UnionId)

	util.Info("%v长链链接成功", req.UnionId)
	// Return OpenId
	c.JSON(http.StatusOK, gin.H{
		"open_id":  req.OpenId,
		"union_id": req.UnionId,
		"room_id":  roomId,
	})

}
