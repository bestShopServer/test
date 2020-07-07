package gocache

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
	"time"
)

var RedisConnPool *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConnPool = &redis.Pool{
		MaxIdle:     config.GetConfig().RedisMaxIdle,
		MaxActive:   config.GetConfig().RedisMaxActive,                                //设MaxActive=0(表示无限大)或者足够大。
		IdleTimeout: time.Duration(config.GetConfig().RedisIdleTimeout) * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭。如果设置成0，空闲连接将不会被关闭。应该设置一个比redis服务端超时时间更短的时间。
		Wait:        true,                                                             //当程序执行get()，无法获得可用连接时，将会暂时阻塞。
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.GetConfig().RedisAddr,
				redis.DialConnectTimeout(time.Duration(config.GetConfig().RedisDialConnectTimeout)*time.Second), //连接Redis超时时间
				redis.DialReadTimeout(time.Duration(config.GetConfig().RedisDialReadTimeout)*time.Second),       //从Redis读取数据超时时间。
				redis.DialWriteTimeout(time.Duration(config.GetConfig().RedisDialWriteTimeout)*time.Second),     //向Redis写入数据超时时间。
			)
			if err != nil {
				return nil, err
			}
			if len(config.GetConfig().RedisAuth) > 0 {
				if _, err := c.Do("AUTH", config.GetConfig().RedisAuth); err != nil {
					defer c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("select", config.GetConfig().RedisDb); err != nil {
				defer c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Close() {
	err := RedisConnPool.Close()
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
}

// Set a key/value
func SetJsonValue(key string, jsonData interface{}) (err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	value, err := json.Marshal(jsonData)
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

// Set a key/value
func Set(key string, value interface{}) error {
	conn := RedisConnPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

//设置redis值, 参数key
func SetKeyExpire(k string, v interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("SET", k, v)
	if err != nil {
		util.Error("set key[%v] error[%v]", k, err.Error())
		return err
	}
	_, err = c.Do("expire", k, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("set key expire[%v] error[%v]", k, err.Error())
		return fmt.Errorf("set expire error", err.Error())
	}
	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConnPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return false
	}

	return exists
}

// Exists check a key
func ConnExists(conn redis.Conn, key string) bool {
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return false
	}
	//util.Debug("key:%v %v", key, exists)

	return exists
}

// Set a key/value
func ConnSetJsonValue(conn redis.Conn, key string, jsonData interface{}) error {
	value, err := json.Marshal(jsonData)
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

	return nil
}

////设置redis值, 参数key,value, expiration(秒)
//func Conn() redis.Conn {
//	c := RedisConnPool.Get()
//	defer c.Close()
//	return c
//}

// Get get a key
func Get(key string) (reply []byte, ok bool, err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	ok = ConnExists(conn, key)
	if ok {
		reply, err = redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return reply, ok, err
		}
	} else {
		util.Error("key:%v不存在", key)
	}

	return reply, ok, nil
}

// Set a key/value
func GetJsonValue(key string, info interface{}) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	bys, ok, err := Get(key)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取信息失败")
	}
	if ok {
		err = json.Unmarshal(bys, &info)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
		}
	}

	return ok, nil
}

func ConnGetJsonValue(conn redis.Conn, key string, info interface{}) (bool, error) {
	ok := ConnExists(conn, key)
	if ok {
		reply, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ok, err
		}
		err = json.Unmarshal(reply, &info)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
		}
	} else {
		util.Error("key:%v不存在", key)
	}

	return ok, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

//设置redis值, 参数key, expiration(秒)
func SetExpire(k string) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	//logging.Debug("keyExpire:%v", setting.RedisSetting.KeyExpire)
	_, err = c.Do("expire", k, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("set key expire[%v] error[%v]", k, err.Error())
		return fmt.Errorf("set expire error", err.Error())
	}
	return nil
}

//设置redis列表值, 参数key,value, expiration(秒)
func RPushKeys(kv ...interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("rpush", kv...)
	if err != nil {
		util.Error("rpush error[%v]", err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	return nil
}

//获取list长度
func GetListLen(k string) (len int64, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	v, err := c.Do("llen", k)
	if err != nil {
		util.Error("llen key[%v] error[%v]", k, err.Error())
		return len, fmt.Errorf("llen error", err.Error())
	}
	util.Info("%v", v)
	//要从byte转string
	len = v.(int64)
	return len, nil
}

//获取list对应索引值
func GetListKey(k string, index int64) (val string, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	v, err := redis.Strings(c.Do("lrange", k, index, index))
	if err != nil {
		util.Error("lrange key[%v] error[%v]", k, err.Error())
		return val, fmt.Errorf("lrange error", err.Error())
	}
	return v[0], nil
}

//删除列表元素
func DelListKey(k string, v interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	bl, err := redis.Bool(c.Do("lrem", k, 1, v))
	if err != nil {
		util.Error("lrem key[%v] error[%v]", k, err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	//logging.Info("删除:%v", bl)
	if !bl {
		util.Error("删除key[%v]值[%v]失败 %v", k, v, bl)
		return fmt.Errorf("删除key[%v]值[%v]失败!", k, v)
	}
	return nil
}

//删除列表元素
func DelListKeyAll(k string, v interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("lrem", k, 0, v)
	if err != nil {
		util.Error("lrem key[%v] error[%v]", k, err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	return nil
}

//设置redis值, 参数key,value, expiration(秒)
func HSetKey(k string, f, v interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("hset", k, f, v)
	if err != nil {
		util.Error("hset key[%v] field[%v] error[%v]", k, f, err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	return nil
}

//设置redis值, 参数key,value, expiration(秒)
func HSetKeys(kfv ...interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("hset", kfv...)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	return nil
}

//读取redis hash值, 参数key,field
func HGetKey(key string, field string) (res []byte, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	res, err = redis.Bytes(c.Do("hget", key, field))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return res, fmt.Errorf("hget error:%v", err.Error())
	}
	return res, err
}

//设置redis值, 参数key,value, expiration(秒)
func HGetAllKeys(key string) (val map[string]string, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	val, err = redis.StringMap(c.Do("hgetall", key))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return val, fmt.Errorf("hgetall error:%v", err.Error())
	}
	//util.Info("hgetall:%v", res)
	return val, err
}

// Exists check a key
func ConnHSetKey(conn redis.Conn, key string, field interface{}, value interface{}) (err error) {
	_, err = conn.Do("hset", key, field, value)
	if err != nil {
		util.Error("hset key[%v] field[%v] error[%v]", key, field, err.Error())
		return fmt.Errorf("set error", err.Error())
	}
	return err
}

//设置redis集合值
func SetSet(k, v interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("sadd", k, v)
	if err != nil {
		util.Error("sadd error[%v]", err.Error())
		return fmt.Errorf("sadd error", err.Error())
	}
	return nil
}

//获取redis集合值所有值
func Smembers(k interface{}) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = c.Do("smembers", k)
	if err != nil {
		util.Error("smembers error[%v]", err.Error())
		return fmt.Errorf("smembers error", err.Error())
	}
	return nil
}

//保存用户所在房间
func SetUserRoom(uid string, roomCode string) (err error) {
	err = SetKeyExpire(uid, roomCode)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	return err
}

//获取用户所在房间
func GetUserRoom(uid string) (roomCode string, err error) {
	by, ok, err := Get(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	if ok {
		roomCode = string(by)
	}

	return roomCode, err
}

//保存房间用户列表
func SetRoomUser(roomCode string, uid string) (err error) {
	err = SetSet("set"+roomCode, uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	err = SetExpire("set" + roomCode)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	return err
}

//添加用户到房间
func AddRoomUser(roomCode string, uid string) (err error) {
	err = SetSet("set"+roomCode, uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	return err
}

//查看房间用户列表
func GetRoomUsers(roomCode string) (users []string, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	users, err = redis.Strings(c.Do("smembers", "set"+roomCode))
	if err != nil {
		util.Error("smembers error[%v]", err.Error())
		return users, fmt.Errorf("获取房间用户失败", err.Error())
	}
	return users, err
}

//判断用户是否在房间内
func CheckUserInRoom(roomCode string, user string) (bl bool, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	bl, err = redis.Bool(c.Do("sismember", "set"+roomCode, user))
	if err != nil {
		util.Error("smembers error[%v]", err.Error())
	}
	return bl, err
}

//删除房间列表中的指定用户
func DelRoomUser(roomCode string, user string) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	_, err = redis.Bool(c.Do("srem", "set"+roomCode, user))
	if err != nil {
		util.Error("srem error[%v]", err.Error())
	}
	return err
}

//删除房间列表中的用户
func DelRoomUsers(roomCode string) (err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	//bl, err := Delete(roomCode)
	err = SetExpire(roomCode)
	if err != nil {
		util.Error("srem error[%v]", err.Error())
	}
	return err
}

//保存房间信息
func SetRoomInfo(roomCode string, info interface{}) (err error) {

	err = SetJsonValue(roomCode, info)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//获取房间信息
func GetRoomInfo(roomCode string, info interface{}) (err error) {
	res, ok, err := Get(roomCode)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取房间信息失败")
		ok = false
	}
	if ok {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(res, &info)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
		}
	}
	return err
}

//判断用户是否在房间内
func CheckRoomExists(roomCode string) (bl bool, err error) {
	c := RedisConnPool.Get()
	defer c.Close()
	bl, err = redis.Bool(c.Do("exists", roomCode))
	if err != nil {
		util.Error("smembers error[%v]", err.Error())
		return false, err
	}
	return bl, err
}

//保存房间投票信息
func SetVoteInfo(roomCode string, votes map[string]bool) (err error) {
	key := "vote" + roomCode
	err = SetJsonValue(key, votes)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间用户投票信息失败")
	}
	return err
}

//获取房间投票信息
func GetVoteInfo(roomCode string) (votes map[string]bool, ok bool, err error) {
	c := RedisConnPool.Get()
	defer c.Close()

	key := "vote" + roomCode
	bys, ok, err := Get(key)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取房间AP点信息失败")
		return votes, ok, err
	}
	if ok {
		err = json.Unmarshal(bys, &votes)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("获取房间用户投票信息失败")
		}
	}
	return votes, ok, err
}

//保存房间用户AP点信息
func SetAPInfo(roomCode string, aps map[string][]int) (err error) {
	//err = HSetKeys("ap"+roomCode, aps)
	err = SetJsonValue("ap"+roomCode, aps)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间用户投票信息失败")
	}
	return err
}

//获取房间投票信息
func GetAPInfo(roomCode string) (apMap map[string][]int, ok bool, err error) {
	c := RedisConnPool.Get()
	defer c.Close()

	key := "ap" + roomCode
	bys, ok, err := Get(key)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取房间AP点信息失败")
		return apMap, ok, err
	}
	if ok {
		err = json.Unmarshal(bys, &apMap)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("获取房间用户投票信息失败")
		}
	}

	return apMap, ok, err
}

//保存房间阶段信息
func SetRoomGame(roomCode string, game interface{}) (err error) {
	err = SetJsonValue("game"+roomCode, game)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//获取房间线索信息
func GetRoomGame(roomCode string, game interface{}) (ok bool, err error) {
	key := "game" + roomCode
	res, ok, err := Get(key)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取房间信息失败")
		return ok, err
	}
	if ok {
		err = json.Unmarshal(res, &game)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间信息失败")
			return ok, err
		}
	}

	return ok, err
}

//保存房间线索信息
func SetRoomClew(roomCode string, clew interface{}) (err error) {
	key := "clew" + roomCode
	err = SetJsonValue(key, clew)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//获取房间线索信息
//func GetRoomClew(roomCode string, clew interface{}) (ok bool, err error) {
func GetRoomClew(roomCode string, res interface{}) (ok bool, err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
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
func SetRoomStage(roomCode string, stage interface{}) (err error) {
	err = SetJsonValue("stage"+roomCode, stage)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("保存房间信息失败")
	}
	return err
}

//获取房间线索信息
func GetRoomStage(roomCode string, stage interface{}) (ok bool, err error) {
	key := "stage" + roomCode
	res, ok, err := Get(key)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("获取房间阶段信息失败")
	}
	if ok {
		err = json.Unmarshal(res, &stage)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			err = errors.New("解析房间阶段信息失败")
			return ok, err
		}
	}
	return ok, err
}

//存储线索
func SetAboutAndTaskAndClewAndExplores(scriptId int, dbparams map[string]interface{}) (err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%v_script_info", scriptId)
	err = SetJsonValue(key, dbparams)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本线索任务失败!")
	}
	return err
}

//存储线索
func GetAboutAndTaskAndClewAndExplores(scriptId int) (dbparams map[string]interface{}, ok bool, err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%v_script_info", scriptId)
	ok, err = GetJsonValue(key, &dbparams)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本线索任务失败!")
	}
	return dbparams, ok, err
}

//保存房间阶段信息
func SetApAndStory(scriptId int, ap []int, story string) (err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
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
	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("设置存储剧本AP点和故事数据过期时间失败!")
	}

	return err
}

func tranApToInts(in string) (out []int, err error) {
	tmp := in[1 : len(in)-1]
	util.Debug("tmp:%v", tmp)
	strs := strings.Fields(tmp)
	for _, str := range strs {
		num, err := strconv.Atoi(str)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return out, err
		}
		out = append(out, num)
	}

	return out, err
}

//获取房间线索信息
func GetApAndStory(scriptId int) (ap []int, story string, ok bool, err error) {
	key := fmt.Sprintf("%vap", scriptId)
	conn := RedisConnPool.Get()
	defer conn.Close()

	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.StringMap(conn.Do("hgetall", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ap, story, ok, err
		}
		//util.Debug("res:%v", reply)
		ap, err = tranApToInts(reply["ap"])
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return ap, story, ok, err
		}
		story = reply["story"]
	} else {
		util.Error("key:%v不存在", key)
	}
	return ap, story, ok, err
}

//记录用户答题得分信息
func SetRoomUserQuestionScore(roomid, userid string, score, maxScore int) (err error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%v_%v_score", roomid, userid)
	err = ConnHSetKey(conn, key, "score", score)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本答题得分数据失败!")
	}
	err = ConnHSetKey(conn, key, "max_score", maxScore)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("存储剧本题目最高得分数据失败!")
	}
	_, err = conn.Do("EXPIRE", key, config.GetConfig().RedisExp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		err = errors.New("设置存储剧本AP点和故事数据过期时间失败!")
	}

	return err
}

//查询用户答题得分信息
func GetRoomUserQuestionScore(roomid, userid string) (score, maxScore int, ok bool, err error) {
	key := fmt.Sprintf("%v_%v_score", roomid, userid)
	conn := RedisConnPool.Get()
	defer conn.Close()

	ok = ConnExists(conn, key)
	if ok {
		reply, err := redis.StringMap(conn.Do("hgetall", key))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			return score, maxScore, ok, err
		}
		score, _ = strconv.Atoi(reply["score"])
		maxScore, _ = strconv.Atoi(reply["max_score"])
	} else {
		util.Error("key:%v不存在", key)
	}
	return score, maxScore, ok, err
}

//获取白名单
func GetWhitelist() ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	us, err := redis.Strings(conn.Do("smembers", "whitelist"))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return us, err
	}
	return us, err
}
