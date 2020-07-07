package point

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"fmt"
	"time"
)

//记录房间数据
func RecordRoomMainData(param model.RoomRecordBase) {
	util.Info("RecordRoomData ...")

	taskRequest := make(map[string]interface{})
	taskRequest["RoomId"] = param.RoomId
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["Owner"] = param.Owner
	taskRequest["UnionId"] = param.UnionId
	taskRequest["Status"] = 0 //1-创建房间
	taskRequest["StartTime"] = time.Now().Format("2006-01-02 15:04:05")
	taskRequest["EndTime"] = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println("RecordRoomData taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.RecordRoomMainBase", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//return playerInfoList, model.ERR_TASK_JSON
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:%v, %v", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("登记数据成功")
		return
	default:
		util.Error("ERROR:%v", dbcode)
		return
	}

}

//记录用户加入房间
func RecordRoomUserJoinData(param model.RoomRecordBase) {
	util.Info("RecordRoomData ...")

	taskRequest := make(map[string]interface{})
	taskRequest["RoomId"] = param.RoomId
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["UnionId"] = param.UnionId
	taskRequest["RoleId"] = param.RoleId
	taskRequest["StartTime"] = time.Now().Format("2006-01-02 15:04:05")
	taskRequest["EndTime"] = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println("RecordRoomData taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.RecordRoomJoinUser", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//return playerInfoList, model.ERR_TASK_JSON
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:%v, %v", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("登记数据成功")
		return
	default:
		util.Error("ERROR:%v", dbcode)
		return
	}
}

//更新房间用户数据
func RoomUserDataUpdate(param model.RoomRecordBase) {
	util.Info("RecordRoomData ...")

	taskRequest := make(map[string]interface{})
	taskRequest["RoomId"] = param.RoomId
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["UnionId"] = param.UnionId
	taskRequest["RoleId"] = param.RoleId
	taskRequest["Score"] = param.Score

	fmt.Println("RecordRoomData taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.RoomUserUpdate", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//return playerInfoList, model.ERR_TASK_JSON
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:%v, %v", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("登记数据成功")
		return
	default:
		util.Error("ERROR:%v", dbcode)
		return
	}

}

func RecordRoomUserExitData(param model.RoomRecordBase) {
	util.Info("RecordRoomData ...")

	taskRequest := make(map[string]interface{})
	taskRequest["RoomId"] = param.RoomId
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["Owner"] = param.Owner
	taskRequest["EndTime"] = time.Now().Format("2006-01-02 15:04:05")
	taskRequest["Status"] = 7 //0创建房间1开始游戏2搜证结束3投票结束5答题结束6评分结束7游戏结束

	fmt.Println("RecordRoomData taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.RoomStatusUpdate", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//return playerInfoList, model.ERR_TASK_JSON
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:%v, %v", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("登记数据成功")
		return
	default:
		util.Error("ERROR:%v", dbcode)
		return
	}
}

func RoomStatusUpdate(param model.RoomRecordBase) {
	util.Info("RecordRoomData ...")

	taskRequest := make(map[string]interface{})
	taskRequest["RoomId"] = param.RoomId
	taskRequest["ScriptId"] = param.ScriptId
	taskRequest["Owner"] = param.Owner
	taskRequest["Status"] = param.Status

	fmt.Println("RecordRoomData taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.RoomStatusUpdate", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//return playerInfoList, model.ERR_TASK_JSON
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:%v, %v", dbcode, len(dbparams))
	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("登记数据成功")
		return
	default:
		util.Error("ERROR:%v", dbcode)
		return
	}

}
