package global

//
//import (
//	"DetectiveMasterServer/config"
//	"DetectiveMasterServer/util"
//	"fmt"
//	"github.com/garyburd/redigo/redis"
//	"time"
//)
//
//var RedisClient *redis.Pool
//
//func init() {
//	// 建立连接池
//	RedisClient = &redis.Pool{
//		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
//		MaxIdle: 16, //最初的连接数量
//		// MaxActive:1000000,    //最大连接数量
//		MaxActive:   0,                 //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
//		IdleTimeout: 300 * time.Second, //连接关闭时间 300秒 （300秒不使用自动关闭）
//		Dial: func() (redis.Conn, error) { //要连接的redis数据库
//			c, err := redis.Dial("tcp", config.GetConfig().RedisAddr)
//			if err != nil {
//				util.Error("TCP ERROR[%v]", err.Error())
//				return nil, err
//			}
//			if len(config.GetConfig().RedisAuth) > 0 {
//				if _, err := c.Do("AUTH", config.GetConfig().RedisAuth); err != nil {
//					c.Close()
//					util.Error("AUTH ERROR[%v]", err.Error())
//					return nil, err
//				}
//			}
//
//			if _, err := c.Do("select", config.GetConfig().RedisDb); err != nil {
//				c.Close()
//				util.Error("SELECT ERROR[%v]", err.Error())
//				return nil, err
//			}
//			return c, nil
//		},
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			if err != nil {
//				util.Error("ERROR[%v]", err.Error())
//				return fmt.Errorf("ping redis error: %s", err)
//			}
//			return nil
//		},
//	}
//}
//
////设置redis值, 参数key,value
//func RedisSetKey(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("SET", k, v)
//	if err != nil {
//		util.Error("set key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
//
////查看redis值, 参数key,value
//func RedisGetKey(k string) (str interface{}, err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	str, err = c.Do("GET", k)
//	if err != nil {
//		util.Error("get key[%v] error[%v]", k, err.Error())
//		return str, fmt.Errorf("set error", err.Error())
//	}
//	return str, nil
//}
//
////删除redis值, 参数key
//func RedisDelKey(k string) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("del", k)
//	if err != nil {
//		util.Error("del key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("del error", err.Error())
//	}
//	return nil
//}
//
////若 key 存在返回 1 ，否则返回 0
//func RedisExistsKey(k string) int {
//	isx := 0
//	c := RedisClient.Get()
//	defer c.Close()
//	is, err := c.Do("exists", k)
//	if err != nil {
//		util.Error("del key[%v] error[%v]", k, err.Error())
//		return isx
//	}
//	isx = is.(int)
//	return isx
//}
//
////设置redis值, 参数key,value, expiration(秒)
//func RedisSetKeyExpire(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	expiration := config.GetConfig().RedisExpSec
//	if expiration > 0 {
//		extime := time.Duration(expiration)
//		_, err = c.Do("SET", k, v, "ex", extime)
//	} else {
//		_, err = c.Do("SET", k, v)
//	}
//	if err != nil {
//		util.Error("set key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
//
////设置redis值, 参数key, expiration(秒)
//func RedisSetExpire(k string) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	extime := time.Duration(config.GetConfig().RedisExpSec)
//	_, err = c.Do("expire", k, extime)
//	if err != nil {
//		util.Error("set key expire[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set expire error", err.Error())
//	}
//	return nil
//}
//
////设置redis值, 参数key,value, expiration(秒)
//func RedisRPushKey(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("rpush", k, v)
//	if err != nil {
//		util.Error("rpush key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
//
////设置redis值, 参数key,value, expiration(秒)
//func RedisDelListKey(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("lrem", k, 0, v)
//	if err != nil {
//		util.Error("lrem key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
//
////设置redis集合值, 参数key,value
//func RedisSetSetKey(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("set", k, v)
//	if err != nil {
//		util.Error("set key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
//
////删除redis集合值, 参数key,value
//func RedisDelSetKey(k string, v interface{}) (err error) {
//	c := RedisClient.Get()
//	defer c.Close()
//	_, err = c.Do("srem", k, v)
//	if err != nil {
//		util.Error("lrem key[%v] error[%v]", k, err.Error())
//		return fmt.Errorf("set error", err.Error())
//	}
//	return nil
//}
