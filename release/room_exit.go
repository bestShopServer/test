package release

//
//// Func: room.exit handler
//func RoomExit(c *gin.Context) {
//	util.Info("RoomExit ...")
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//
//	// Get Page and Scene and Width From Param
//	var req model.RoomExitReq
//	err := c.Bind(&req)
//
//	// Optional Fields List
//	var optionalFields []string
//
//	// Check Param
//	if !CheckParams(req, "RoomExit", err, optionalFields) {
//		util.Error("ERR_WRONG_FORMAT")
//		ResponseErr(model.ERR_WRONG_FORMAT, c)
//		return
//	}
//
//	// Find RoomId In Room Bucket
//	//uid := req.OpenId
//	uid := req.UnionId
//	rid := req.RoomId
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
//	roomInfo := model.RoomInfo{}
//	de := json.Unmarshal(b, &roomInfo)
//	if de != nil {
//		//util.Logger(util.ERROR_LEVEL, "RoomExit", "Decoding Room Info Err:"+de.Error())
//		util.Error("RoomExit Decoding Room Info ERROR[%v]", de.Error())
//	}
//
//	// If You Are In This Room Then Remove Your OpenId
//	ois := roomInfo.UnionIdSlice
//	exist := false
//	for k, v := range roomInfo.UnionIdSlice {
//		if v == uid {
//			exist = true
//			roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice[:k], roomInfo.UnionIdSlice[k+1:]...)
//			break
//		}
//	}
//	if !exist {
//		util.Error("ERR_BELONG")
//		ResponseErr(model.ERR_BELONG, c)
//		return
//	}
//
//	// If You Are In Vote Stage Then Remove Your OpenId And Send Game Vote Message
//	vb := boltdb.View([]byte(rid), "VoteBucket")
//	if vb != nil {
//		var votes map[string]bool
//		json.Unmarshal(vb, &votes)
//		if votes[uid] == false {
//			votes[uid] = true
//			notVoteNum := global.CalcNotVoteNum(votes)
//			encodingVotes, _ := json.Marshal(votes)
//			boltdb.CreateOrUpdate([]byte(rid), encodingVotes, "VoteBucket")
//			go websocket.SendGameVoteMsg(roomInfo.UnionIdSlice, notVoteNum)
//		}
//	}
//
//	// Update BoltDB
//	roomInfoEncoding, err := json.Marshal(roomInfo)
//	if err != nil {
//		//util.Logger(util.ERROR_LEVEL, "RoomExit", "Room Info Encoding Err:"+err.Error())
//		util.Error("RoomExit Room Info Encoding ERROR[%v]", err.Error())
//	}
//	boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
//
//	// Update User Cache
//	global.DeleteUserCache(uid, rid)
//
//	// Send Room Exit Msg
//	go websocket.SendRoomExitMessage(ois, uid)
//
//	// Return
//	resp := model.RoomExitResp{}
//	resp.Params = "success"
//	c.JSON(http.StatusOK, resp)
//}
