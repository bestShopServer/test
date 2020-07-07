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

func PlayerSliceModify(s *[]model.PlayerInfo, index int, value model.PlayerInfo) {
	rear := append([]model.PlayerInfo{}, (*s)[index+1:]...)
	*s = append(append((*s)[:index], value), rear...)
}

func UserInfoSliceModify(s *[]model.UserInfo, index int, value model.UserInfo) {
	rear := append([]model.UserInfo{}, (*s)[index+1:]...)
	*s = append(append((*s)[:index], value), rear...)
}

func GameVote(c *gin.Context) {

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get openId & roomId From Param
	var req model.GameVoteReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "GameVote", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	rid := req.RoomId
	//uid := req.OpenId
	uid := req.UnionId

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

	// Find Room in RoomBucket
	//b := boltdb.View([]byte(rid), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	//查看用户是否在白名单
	isWhite := false
	whiteUsers, err := gocache.GetWhitelist()
	if err != nil {
		util.Error("获取白名单失败, ERROR:%v", err.Error())
	}
	util.Info("白名单数据:%+v", whiteUsers)
	for _, tmp := range whiteUsers {
		if tmp == uid {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameVote", "Decoding Room Info Err:"+de.Error())
	//	util.Error("GameVote Decoding Room Info ERROR[%v", de.Error())
	//}
	//读取房间信息
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
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	//查询被投用户的union_id
	var mUnionId string
	for i, player := range roomInfo.PlayerSlice {
		if player.Role.Id == req.RoleId {
			mUnionId = player.UnionId
			if player.Role.Murderer {
				player.VoteResult = true //若是凶手则投凶正确
			}
			PlayerSliceModify(&roomInfo.PlayerSlice, i, player)
			break
		}
	}
	util.Debug("投的凶手ID:%v", mUnionId)

	//处理被投数据
	for i, us := range roomInfo.UserInfoSlice {
		if us.UnionId == mUnionId {
			us.CoverVoteNum += 1 //增加被投次数
			UserInfoSliceModify(&roomInfo.UserInfoSlice, i, us)
		}
		if us.UnionId == req.UnionId {
			us.VoteUser = mUnionId //记录投凶手票
			UserInfoSliceModify(&roomInfo.UserInfoSlice, i, us)
		}
	}
	util.Debug("更新房间数据...")
	//roomInfoEncoding, err := json.Marshal(roomInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "RoomNew", "Room Info Encoding Err:"+err.Error())
	//	util.Error("RoomNew Room Info Encoding ERROR[%v]", err.Error())
	//}
	//
	//boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
	err = gocache.SetRoomInfo(rid, roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Debug("处理投票记录...")

	// Find Room In VoteBucket
	//vb := boltdb.View([]byte(rid), "VoteBucket")
	//if vb == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Vote Map
	//var votes map[string]bool
	//err = json.Unmarshal(vb, &votes)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameVote", "Decoding Votes Err:"+err.Error())
	//	util.Error("GameVote Decoding Votes ERROR[%v]", err.Error())
	//}
	votes, bl, err := gocache.GetVoteInfo(rid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_NOT_ALL_VOTED, c)
		return
	}
	if !bl {
		util.Error("ERR_ROOM_NOT_EXIST")
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	if votes[uid] == true {
		util.Error("ERR_HAS_VOTED")
		ResponseErr(model.ERR_HAS_VOTED, c)
		return
	}
	util.Debug("GetVoteInfo...")

	votes[uid] = true
	notVoteNum := global.CalcNotVoteNum(votes)

	if isWhite { //处理白名单数据
		util.Debug("白名单...")
		for i, _ := range votes {
			votes[i] = true //自动全部投票
		}
	}
	// Update BoltDB
	//encodingVotes, err := json.Marshal(votes)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameVote", "Encoding Votes Err:"+err.Error())
	//	util.Error("GameVote Encoding Votes ERROR[%v]", err.Error())
	//}
	//boltdb.CreateOrUpdate([]byte(rid), encodingVotes, "VoteBucket")
	err = gocache.SetVoteInfo(rid, votes)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_NOT_ALL_VOTED, c)
		return
	}
	util.Debug("SetVoteInfo...")

	// Send Game Vote Msg
	if isWhite {
		util.Debug("白名单...")
		go websocket.SendGameVoteMsg(roomInfo.UnionIdSlice, 0)
	} else {
		go websocket.SendGameVoteMsg(roomInfo.UnionIdSlice, notVoteNum)
	}

	// Return
	resp := model.GameVoteResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, &resp)
}
