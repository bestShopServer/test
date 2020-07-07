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

// Func: wx.exit handler
func WxExit(c *gin.Context) {
	util.Info("WxExit ...")

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get Params
	var req model.WxExitReq
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

	////Get RoomId
	//rid, ok := global.UserCache[uid]
	//if !ok {
	//	ResponseErr(model.ERR_NOT_IN_ROOM, c)
	//	return
	//}
	//获取用户房间号
	rid, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_NOT_IN_ROOM, c)
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
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "RoomExit", "Decoding Room Info Err:"+de.Error())
	//	util.Error("RoomExit Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(rid, &roomInfo)
	if err != nil {
		util.Error(" ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You Are In This Room Then Remove Your OpenId
	ois := roomInfo.UnionIdSlice
	exist := false
	for k, v := range roomInfo.UnionIdSlice {
		if v == uid {
			exist = true
			roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice[:k], roomInfo.UnionIdSlice[k+1:]...)
			break
		}
	}
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	//删除房间用户信息
	for i, u := range roomInfo.UserInfoSlice {
		if u.UnionId == uid {
			roomInfo.UserInfoSlice = append(roomInfo.UserInfoSlice[:i], roomInfo.UserInfoSlice[i+1:]...)
			break
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
	//		//encodingVotes, _ := json.Marshal(votes)
	//		//boltdb.CreateOrUpdate([]byte(rid), encodingVotes, "VoteBucket")
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
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Update User Cache
	//global.DeleteUserCache(uid, rid)
	err = gocache.DelRoomUser(rid, uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Send Room Exit Msg
	go websocket.SendRoomExitMessage(ois, uid)

	// Return
	resp := model.WxExitResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, resp)
}
