package gocache

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

func ConnGetUserRoom(conn redis.Conn, uid string) (roomCode string, err error) {
	ok := ConnExists(conn, uid)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", uid))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return roomCode, err
		}
		roomCode = string(reply)
	} else {
		util.Error("key:%v不存在", uid)
		err = fmt.Errorf("%v查询不存在", uid)
	}
	return roomCode, err
}

//判断用户是否在房间内
func ConnCheckRoomExists(conn redis.Conn, roomCode string) (bl bool, err error) {
	bl, err = redis.Bool(conn.Do("exists", roomCode))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return false, err
	}
	return bl, err
}

//获取房间信息
func ConnGetRoomInfo(conn redis.Conn, roomCode string, info interface{}) (err error) {
	ok := ConnExists(conn, roomCode)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", roomCode))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return err
		}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(reply, &info)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
		}
	} else {
		util.Error("key:%v不存在", roomCode)
	}
	return err
}

//保存房间信息
func ConnSetRoomInfo(conn redis.Conn, roomCode string, info interface{}) (err error) {
	value, err := json.Marshal(info)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	//util.Debug("conn值没了:%+v", string(value))
	_, err = conn.Do("SET", roomCode, value)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("EXPIRE", roomCode, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	return err
}

//存储线索
func ConnGetAboutAndTaskAndClewAndExplores(conn redis.Conn, scriptId int) (dbparams map[string]interface{}, ok bool, err error) {
	key := fmt.Sprintf("%v_script_info", scriptId)
	ok, err = ConnGetJsonValue(conn, key, &dbparams)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return dbparams, ok, err
	}
	return dbparams, ok, err
}

//存储线索
func ConnSetAboutAndTaskAndClewAndExplores(conn redis.Conn, scriptId int, dbparams map[string]interface{}) (err error) {
	key := fmt.Sprintf("%v_script_info", scriptId)
	value, err := json.Marshal(dbparams)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp*10)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	return err
}

//保存房间阶段信息
func ConnSetApAndStory(conn redis.Conn, scriptId int, ap []int, story, final string) (err error) {
	key := fmt.Sprintf("%vap", scriptId)
	err = ConnHSetKey(conn, key, "ap", ap)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本AP点和故事数据失败!")
	}
	err = ConnHSetKey(conn, key, "story", story)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本AP点和故事数据失败!")
	}
	err = ConnHSetKey(conn, key, "final", final)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本AP点和故事数据失败!")
	}
	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp*10)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("设置存储剧本AP点和故事数据过期时间失败!")
	}

	return err
}

//获取房间线索信息
func ConnGetApAndStory(conn redis.Conn, scriptId int) (ap []int, story, final string, ok bool, err error) {
	key := fmt.Sprintf("%vap", scriptId)
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.StringMap(conn.Do("hgetall", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ap, story, final, ok, err
		}
		//util.Debug("res:%v", reply)
		ap, err = tranApToInts(reply["ap"])
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ap, story, final, ok, err
		}
		story = reply["story"]
		final = reply["final"]
	} else {
		util.Error("key:%v不存在", key)
	}

	return ap, story, final, ok, err
}

//获取房间线索信息
func ConnGetRoomGame(conn redis.Conn, roomCode string, game interface{}) (ok bool, err error) {
	key := "game" + roomCode
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
		err = json.Unmarshal(reply, &game)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
			return ok, err
		}
	} else {
		util.Error("key:%v不存在", key)
	}

	return ok, err
}

//保存房间阶段信息
func ConnSetRoomGame(conn redis.Conn, roomCode string, game interface{}) (err error) {
	key := "game" + roomCode
	err = ConnSetJsonValue(conn, key, game)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储房间信息失败")
		return err
	}
	return err
}

//获取房间线索信息
func ConnGetRoomStage(conn redis.Conn, roomCode string, stage interface{}) (ok bool, err error) {
	key := "stage" + roomCode
	ok = ConnExists(conn, key)
	if ok {

		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
		err = json.Unmarshal(reply, &stage)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间阶段信息失败")
			return ok, err

		}
	} else {
		util.Error("key:%v不存在", key)
	}
	return ok, err
}

//获取房间投票信息
func ConnGetAPInfo(conn redis.Conn, roomCode string) (apMap map[string][]int, ok bool, err error) {
	key := "ap" + roomCode
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return apMap, ok, err
		}
		err = json.Unmarshal(reply, &apMap)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("获取房间用户投票信息失败")
		}
	} else {
		util.Error("key:%v不存在", key)
	}

	return apMap, ok, err
}

//保存房间用户AP点信息
func ConnSetAPInfo(conn redis.Conn, roomCode string, aps map[string][]int) (err error) {
	key := "ap" + roomCode
	value, err := json.Marshal(aps)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	return err
}

//获取房间投票信息
func ConnGetVoteInfo(conn redis.Conn, roomCode string) (votes map[string]bool, ok bool, err error) {
	key := "vote" + roomCode
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return votes, ok, err
		}
		err = json.Unmarshal(reply, &votes)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("获取房间用户投票信息失败")
		}
	} else {
		util.Error("key:%v不存在", key)
	}

	return votes, ok, err
}

//保存房间投票信息
func ConnSetVoteInfo(conn redis.Conn, roomCode string, votes map[string]bool) (err error) {
	key := "vote" + roomCode
	value, err := json.Marshal(votes)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	return err
}

//保存房间线索信息
func ConnSetRoomClew(conn redis.Conn, roomCode string, clew interface{}) (err error) {
	key := "clew" + roomCode
	err = ConnSetJsonValue(conn, key, clew)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//获取房间线索信息
func ConnGetRoomClew(conn redis.Conn, roomCode string, res interface{}) (ok bool, err error) {
	key := "clew" + roomCode
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
		err = json.Unmarshal(reply, &res)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
	} else {
		util.Error("key:%v不存在", key)
	}
	return ok, err
}

//保存房间阶段信息
func ConnSetRoomStage(conn redis.Conn, roomCode string, stage interface{}) (err error) {
	key := "stage" + roomCode
	err = ConnSetJsonValue(conn, key, stage)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//保存剧本人物信息
func ConnSetScriptPeopleInfo(conn redis.Conn, scriptId, roleId int, info interface{}) (err error) {
	key := fmt.Sprintf("%v_%v_role", scriptId, roleId)
	err = ConnSetJsonValue(conn, key, info)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存角色数据信息失败")
	}
	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp*10)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("设置存储剧本AP点和故事数据过期时间失败!")
	}
	return err
}

//获取房间线索信息
func ConnGetScriptPeopleInfo(conn redis.Conn, scriptId, roleId int, info interface{}) (ok bool, err error) {
	key := fmt.Sprintf("%v_%v_role", scriptId, roleId)
	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
		err = json.Unmarshal(reply, &info)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析剧本角色信息失败")
			return ok, err
		}
	} else {
		util.Error("key:%v不存在", key)
	}
	return ok, err
}

//保存房间用户列表
func ConnSetRoomUser(conn redis.Conn, roomCode string, uid string) (err error) {
	key := "set" + roomCode
	//err = SetSet("set"+roomCode, uid)
	_, err = conn.Do("sadd", key, uid)
	if err != nil {
		util.Error("sadd error[%v]", err.Error())
		return fmt.Errorf("sadd error", err.Error())
	}
	//设置过期时间
	_, err = conn.Do("expire", key, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("set key expire[%v] error[%v]", roomCode, err.Error())
		return fmt.Errorf("set expire error", err.Error())
	}
	return err
}
