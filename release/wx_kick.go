package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: wx.kick handler
func WxKick(c *gin.Context) {
	util.Info("WxKick ...")

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get Params
	var req model.WxKickReq
	var ois []string

	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Params
	if !CheckParams(req, "WxExit", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	// Get OpenId
	//uid := req.OpenId
	uid := req.UnionId

	//// Get RoomId
	//rid, ok := global.UserCache[uid]
	//if !ok {
	//	util.Info("ERR_NOT_IN_ROOM")
	//	ResponseErr(model.ERR_NOT_IN_ROOM, c)
	//	return
	//}
	rid, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}

	// Get Kick User RoomId
	//krid, ok := global.UserCache[req.KickUnionId]
	//if !ok {
	//	util.Info("ERR_NOT_IN_ROOM")
	//	ResponseErr(model.ERR_NOT_IN_ROOM, c)
	//	return
	//}

	krid, err := gocache.GetUserRoom(req.KickUnionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}

	util.Info("房间号:%v 被踢着房间号:%v", rid, krid)
	//判断是否为同一个房间
	if krid != rid {
		util.Info("ERR_NOT_IN_ROOM")
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
		return
	}

	//// Find KV in Room Bucket
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
	//	//util.Logger(util.ERROR_LEVEL, "RoomExit", "Decoding Room Info Err:"+de.Error())
	//	util.Error("RoomExit Decoding Room Info ERROR[%v]", de.Error())
	//}
	err = gocache.GetRoomInfo(rid, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	//判断是否是房主权限踢人
	util.Info("房主[%v]", roomInfo.Owner)
	//if uid != roomInfo.UnionIdSlice[0] {
	if uid != roomInfo.Owner {
		util.Error("非房主无权操作")
		ResponseErr(model.ERR_NOT_ROOM_OWNER, c)
		return
	}
	//if req.KickUnionId == roomInfo.UnionIdSlice[0] {
	if req.KickUnionId == roomInfo.Owner {
		util.Error("不能踢出房主")
		ResponseErr(model.ERR_NOT_KICK_ROOM_OWNER, c)
		return
	}

	// If You Are In This Room Then Remove Your OpenId
	//ois := roomInfo.UnionIdSlice
	ois = append(ois, roomInfo.UnionIdSlice[:]...)
	util.Info("房间用户:%v", ois)
	exist := false
	for k, v := range roomInfo.UnionIdSlice {
		if v == req.KickUnionId {
			exist = true
			roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice[:k], roomInfo.UnionIdSlice[k+1:]...)
			break
		}
	}
	//util.Info("房间用户:%v", ois)
	//util.Info("房间剩余用户:%v", roomInfo.UnionIdSlice)
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	//删除被踢用户信息
	for i, u := range roomInfo.UserInfoSlice {
		if u.UnionId == req.KickUnionId {
			util.Info("删除房间[%v]用户[%v]", rid, req.KickUnionId)
			roomInfo.UserInfoSlice = append(roomInfo.UserInfoSlice[:i], roomInfo.UserInfoSlice[i+1:]...)
		}
	}

	// If You Are In Vote Stage Then Remove Your OpenId And Send Game Vote Message
	//vb := boltdb.View([]byte(rid), "VoteBucket")
	//if vb != nil {
	//	var votes map[string]bool
	//	json.Unmarshal(vb, &votes)
	//	if votes[uid] == false {
	//		votes[uid] = true
	//		notVoteNum := global.CalcNotVoteNum(votes)
	//		encodingVotes, _ := json.Marshal(votes)
	//		boltdb.CreateOrUpdate([]byte(rid), encodingVotes, "VoteBucket")
	//		go websocket.SendGameVoteMsg(roomInfo.UnionIdSlice, notVoteNum)
	//	}
	//}
	votes, bl, err := gocache.GetVoteInfo(rid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_HAS_VOTED, c)
		return
	}
	if bl {
		if votes[uid] == false {
			votes[uid] = true
			notVoteNum := global.CalcNotVoteNum(votes)
			err = gocache.SetVoteInfo(rid, votes)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_HAS_VOTED, c)
				return
			}
			go websocket.SendGameVoteMsg(roomInfo.UnionIdSlice, notVoteNum)
		}
	}

	//// Update BoltDB
	//roomInfoEncoding, err := json.Marshal(roomInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "RoomExit", "Room Info Encoding Err:"+err.Error())
	//	util.Error("RoomExit Room Info Encoding ERROR[%v]", err.Error())
	//}
	//boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
	err = gocache.SetRoomInfo(rid, roomInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ENTERED_ROOM, c)
		return
	}

	// Update User Cache
	//global.DeleteUserCache(req.KickUnionId, rid)
	err = gocache.DelRoomUser(rid, req.KickUnionId)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ENTERED_ROOM, c)
		return
	}

	// Send Room Exit Msg
	util.Info("通知用户:%v", ois)
	go websocket.SendRoomExitMessage(ois, req.KickUnionId)
	util.Info("踢出成功!")

	// Return
	resp := model.WxExitResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, resp)
}
