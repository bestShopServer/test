package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"net/http"
	"strconv"
)

// Func: wx.join handler
func WxJoin(c *gin.Context) {
	util.Info("WxJoin ...")
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var user model.UserInfo
	// Get openId & roomId From Param
	var req model.WxJoinReq
	roomRoleNum := 0 //房间可被选人数

	err := c.Bind(&req)

	// Optional Fields List
	//var optionalFields []string
	optionalFields := []string{"UserName"}

	// Check Param
	if !CheckParams(req, "WxJoin", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId
	user.UnionId = req.UnionId
	user.Name = req.UserName
	if len(user.Name) == 0 {
		user.Name = "正在进入"
	}

	// Find Room In Room Cache
	//_, ok := global.RoomCache[rid]
	//if !ok {
	//	util.Error("ERR_ROOM_NOT_EXIST [%v]", rid)
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	ok, err := gocache.CheckRoomExists(rid)
	if err != nil {
		util.Error("ERR_ROOM_NOT_EXIST [%v]", rid)
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if !ok {
		util.Error("ERR_ROOM_NOT_EXIST [%v]", rid)
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

	//房间游戏已结束，则无法再加入房间
	if rif.Status == 2 { //标识房间已结束
		util.Error("房间:%v 状态:%v", rid, rif.Status)
		ResponseErr(model.ERR_ROOM_ALREADY_OVER, c)
		return
	}

	roomCode, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Info("用户:%v在房间:%v", uid, roomCode)
	if roomCode != "" && roomCode != rid {
		//判断之前的房间是否已结束
		rif2 := model.RoomInfo{}
		ok, err := gocache.CheckRoomExists(roomCode)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}
		if ok {
			err = gocache.GetRoomInfo(roomCode, &rif2)
			if err != nil {
				util.Error("ERROR[%v]", err.Error())
				ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
				return
			}
			if rif2.Status != 2 { //标识房间已结束
				util.Error("房间:%v 状态:%v", roomCode, rif2.Status)
				//ResponseErr(model.ERR_ROOM_NOT_OVER, c)
				msg := fmt.Sprintf("上一房间%v游戏未结束，不允许加入新房间!", roomCode)
				ResponseNotice(c, model.ERR_ROOM_NOT_OVER, msg)
				return
			}
		}
	}

	//判断用户是否在房间
	bl, err := gocache.CheckUserInRoom(rid, uid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if bl {
		util.Info("用户:%v已在房间:%v", uid, rid)
	} else {
		//重新计算房间可被选人
		for _, tmp := range rif.PlayerSlice {
			if tmp.Role.Choice {
				roomRoleNum += 1
			}
		}
		util.Info("剧本可被选角色数:%d, 当前房间人数:%v", roomRoleNum, len(rif.UnionIdSlice))

		// If openIdSlice length eq playerSlice
		//if len(rif.UnionIdSlice) == len(rif.PlayerSlice) {
		if len(rif.UnionIdSlice) >= roomRoleNum {
			util.Error("ERR_ROOM_IS_FULL")
			ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}

		//用户超过房间人数也需要报错
		if len(rif.UnionIdSlice) > rif.Num {
			util.Error("ERR_ROOM_IS_FULL")
			ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}

		// If openIdSlice Contains openId
		exist := false
		for _, v := range rif.UnionIdSlice {
			if v == uid {
				exist = true
				break
			}
		}

		useExit := false
		for _, u := range rif.UserInfoSlice {
			if u.UnionId == uid {
				useExit = true
				break
			}
		}

		//没数据则增加用户信息
		if !useExit {
			util.Info("剧本价格:%f", rif.Price)
			if rif.Price > 0.001 {
				var cost model.ScriptUserCostReq
				cost.ScriptId = rif.ScriptId
				cost.UnionIds = append(cost.UnionIds, req.UnionId)

				//查询付款
				us, code := CheckUserScriptCost(cost)
				util.Debug("查询用户剧本付费情况,code:%v", code)
				if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
					util.Error("Get User ScriptId Err [%v]", code)
					ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
					return
				}
				if code == global.ERR_DB_NOTFOUND_DATA { //未付款
					user.IsPay = 1
				} else {
					user.IsPay = us[0].CostFlag
				}
			}

			rif.UserInfoSlice = append(rif.UserInfoSlice, user) //添加用户
		}

		if !exist {
			// Add openId to openIdSlice
			rif.UnionIdSlice = append(rif.UnionIdSlice, uid)

			//// Update BoltDB
			//roomInfoEncoding, err := json.Marshal(rif)
			//if err != nil {
			//	util.Logger(util.ERROR_LEVEL, "RoomEnter", "Room Info Encoding Err:"+err.Error())
			//}
			//boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
			err = gocache.SetRoomInfo(rid, rif)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_ENTERED_ROOM, c)
				return
			}

			//// Update UserCache & RoomCache
			//global.SetUserCache(uid, rid)
			//global.AddUserToRoomCache(uid, rid)
			err = gocache.SetUserRoom(uid, rid) //保存用户房间号
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
				return
			}
			err = gocache.AddRoomUser(rid, uid) //保存房间用户列表
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
				return
			}

			// Send Room Enter Msg
			//go websocket.SendRoomEnterMessage(rif.UnionIdSlice, uid)
			go websocket.SendRoomEnterMessage(rif.UnionIdSlice, user)

		}
	}

	// Return Room Info
	count := CountUnchoosedPeople(rif.PlayerSlice)

	resp := model.WxJoinResp{}
	//resp.Params.Path = "/pages/loadingnew/loadingnew?is_create=0&script_id=" +
	//	strconv.Itoa(rif.ScriptId) + "&unchoosed_num=" + strconv.Itoa(count) +
	//	"&open_id=" + uid + "&union_id=" + uid + "&room_id=" + rid
	resp.Params.Path = "/pages/pick/pick?is_create=0&script_id=" +
		strconv.Itoa(rif.ScriptId) + "&unchoosed_num=" + strconv.Itoa(count) +
		"&open_id=" + uid + "&union_id=" + uid + "&room_id=" + rid

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err = jsonEncoder.Encode(resp); err != nil {
		ResponseErr(model.ERR_DEFAULT, c)
		return
	}

	c.Data(http.StatusOK, "application/json", bf.Bytes())
}

// Func: Calc Unchoosed People
func CountUnchoosedPeople(sli []model.PlayerInfo) int {
	var count int
	for _, v := range sli {
		if v.UnionId == "" {
			count++
		}
	}
	return count
}
