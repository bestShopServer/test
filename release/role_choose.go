package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Func: role.choose handler
func RoleChoose(c *gin.Context) {
	util.Info("RoleChoose ...")
	//defer global.RoomInfoMutex.Unlock()
	//global.RoomInfoMutex.Lock()

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get openId & roomId From Param
	var req model.RoleChooseReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "RoleChoose", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	roomId := req.RoomId
	roleId := req.RoleId
	//userId := req.OpenId
	userId := req.UnionId

	// Find Room In Room Cache
	//_, ok := global.RoomCache[roomId]
	//if !ok {
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}
	ok, err := gocache.CheckRoomExists(roomId)
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

	// Find RoomId in Room Bucket
	//b := boltdb.View([]byte(roomId), "RoomBucket")
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
		if tmp == userId {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//de := json.Unmarshal(b, &roomInfo)
	//if de != nil {
	//	//util.Logger(util.ERROR_LEVEL, "RoleChoose", "Decoding Room Info Err:"+de.Error())
	//	util.Error("RoleChoose Decoding Room Info ERROR[%v]", de.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If openId not Exist in Room
	exist := false
	for _, v := range roomInfo.UnionIdSlice {
		if v == userId {
			exist = true
			break
		}
	}
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	// If You has choose a Role
	for _, v := range roomInfo.PlayerSlice {
		if v.UnionId == userId {
			util.Error("ERR_ROLE_SELECT")
			ResponseErr(model.ERR_ROLE_SELECT, c)
			return
		}
	}

	// Find RoleId in PlayerSlice
	update := false
	for i, v := range roomInfo.PlayerSlice {
		if v.Role.Id == roleId {
			if v.UnionId != "" {
				util.Error("ERR_ROLE_SELECTED")
				ResponseErr(model.ERR_ROLE_SELECTED, c)
				return
			} else {
				roomInfo.PlayerSlice[i].UnionId = userId
				update = true
				break
			}
		}
	}

	// Update Key Value in DB
	if update {
		//若是白名单用户，自动选择其他角色
		if isWhite {
			for i, v := range roomInfo.PlayerSlice {
				if len(v.UnionId) == 0 {
					autoUid := strconv.FormatInt(time.Now().UnixNano(), 10)
					roomInfo.PlayerSlice[i].UnionId = autoUid
					roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice, autoUid)
					user := model.UserInfo{}
					user.UnionId = autoUid
					user.Name = fmt.Sprintf("自动%v", i)

					//发送长链
					go websocket.SendRoleChooseMessage(roomInfo.UnionIdSlice, autoUid, v.Role.Id)
				}
			}
		}

		//roomInfoEncoding, err := json.Marshal(roomInfo)
		//if err != nil {
		//	//util.Logger(util.ERROR_LEVEL, "RoleChoose", "Room Info Encoding Err:"+err.Error())
		//	util.Error("RoleChoose Room Info Encoding ERROR[%v]", err.Error())
		//}
		//boltdb.CreateOrUpdate([]byte(roomId), roomInfoEncoding, "RoomBucket")
		err = gocache.SetRoomInfo(roomId, roomInfo)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ENTERED_ROOM, c)
			return
		}
	}

	// Send Broadcast Msg
	if update {
		go websocket.SendRoleChooseMessage(roomInfo.UnionIdSlice, userId, roleId)

		//更新用户角色
		param := model.RoomRecordBase{}
		param.ScriptId = roomInfo.ScriptId
		param.RoomId = roomId
		param.UnionId = userId
		param.RoleId = roleId
		go point.RoomUserDataUpdate(param)
	}
	util.Info("success")

	// Return Success
	resp := model.RoleChooseResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, resp)
}
