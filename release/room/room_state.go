package room

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
)

//房间状态
func ExitLastRoom(uid string) (err error) {
	util.Info("ExitLastRoom ...")
	noticeUnionIds := []string{}
	//判断用户是否在之前房间，若存在则释放之前的房间
	roomCodeOld, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	//用户存在上一房间
	if len(roomCodeOld) > 0 {
		roomInfo := model.RoomInfo{}
		//查询房间数据
		err = gocache.GetRoomInfo(roomCodeOld, &roomInfo)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			return err
		}

		isChoosePlayer := false
		for _, v := range roomInfo.PlayerSlice {
			if v.UnionId == uid {
				isChoosePlayer = true
				break
			}
		}

		if isChoosePlayer {
			//删除记录用户房间信息
			rooms, err := gocache.GetRoomUsers(roomCodeOld)
			for _, uid := range rooms {
				roomCode, err := gocache.GetUserRoom(uid)
				if err != nil {
					util.Error("ERROR[%v]", err.Error())
					return err
				}
				//用户在当前房间的删除记录
				if roomCode == roomCodeOld {
					noticeUnionIds = append(noticeUnionIds, uid)
					_, err = gocache.Delete(uid)
					if err != nil {
						util.Error("ERROR[%v]", err.Error())
						return err
					}
				}
			}

			//删除房间信息
			_, err = gocache.Delete(roomCodeOld)
			if err != nil {
				util.Error("ERROR[%v]", err.Error())
				return err
			}

			//记录用户离开房间
			param := model.RoomRecordBase{}
			param.ScriptId = roomInfo.ScriptId
			param.RoomId = roomCodeOld
			param.Owner = roomInfo.Owner
			param.UnionId = uid
			go point.RecordRoomUserExitData(param)

			//for i, v := range roomInfo.UnionIdSlice {
			//	if v == uid {
			//		roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice[:i], roomInfo.UnionIdSlice[i+1:]...)
			//	}
			//}
			// util.Debug("通知用户:%+v", roomInfo.UnionIdSlice)
			for i, v := range noticeUnionIds {
				if v == uid {
					noticeUnionIds = append(noticeUnionIds[:i], noticeUnionIds[i+1:]...)
				}
			}
			util.Debug("通知用户:%+v", noticeUnionIds)

			//长链通知用户房间已删除
			go websocket.SendRoomDeleteMessage(noticeUnionIds)

		}
		//长链通知其他用户，自己退出房间
		go websocket.SendRoomExitMessage(noticeUnionIds, uid)
	}

	return err
}
