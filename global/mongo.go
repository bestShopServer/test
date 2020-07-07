package global

//
//import (
//	"DetectiveMasterServer/config"
//	"DetectiveMasterServer/util"
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	"time"
//)
//
//var MongoDB *mgo.Database
//
//func InitMongoDB() {
//	dialInfo := mgo.DialInfo{
//		Addrs:     []string{},
//		Direct:    false,
//		Timeout:   time.Second * 1,
//		PoolLimit: 4096,
//	}
//	//session, err := mgo.Dial(config.GetConfig().MongoAddr)
//	session, err := mgo.DialWithInfo(&dialInfo)
//	if err != nil {
//		util.Error("链接MongoDB失败[%v]", err.Error())
//	}
//	defer session.Close()
//	util.Info("链接MongoDB成功...")
//	session.SetMode(mgo.Monotonic, true)
//
//	MongoDB = session.DB(config.GetConfig().MongoDb)
//	if len(config.GetConfig().MongoUser) > 0 {
//		err = MongoDB.Login(config.GetConfig().MongoUser, config.GetConfig().MongoPasswd)
//		if err != nil {
//			util.Error("链接MongoDB失败[%v]", err.Error())
//		}
//	}
//	util.Info("链接成功...")
//
//}
//
//var CollectionRoom = "room"
//var CollectionRoomCache = "roomIds"
//var CollectionUserCache = "userIds"
//
//
////创建房间信息表，用于存放房间内人员信息
//func CreateRoom() error {
//	err := MongoDB.C(CollectionRoom).EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true})
//	if err != nil {
//		util.Error("链接房间唯一索引失败[%v]", err.Error())
//		return err
//	}
//	return err
//}
//
////登记房间用户
//func InsertRoomUser(roomId, unionId string) error {
//	err := MongoDB.C(CollectionRoom).Insert(bson.M{"id": roomId, "user": unionId, "time": time.Now().Format("2006-01-02 15:04:05")})
//	if err != nil {
//		util.Error("链接房间唯一索引失败[%v]", err.Error())
//		return err
//	}
//	return nil
//}
//
////查询房间的所有者
//func QueryRoomOwner(roomId string) (unionId string, err error) {
//
//	return unionId, err
//}
//
//
////创建房间并设置房间过期
////格式:{"id":"1212", "owner":"wxUnionId", "time":new Date()}
//func CreateRoomCache() error {
//	//创建房间号唯一索引
//	err := MongoDB.C(CollectionRoomCache).EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true})
//	if err != nil {
//		util.Error("房间唯一索引失败[%v]", err.Error())
//		return err
//	}
//	//设置过期时间
//	err = MongoDB.C("roomIds").EnsureIndex(mgo.Index{Key: []string{"time"},
//		ExpireAfter: 24 * time.Hour}) //24小时房间过期
//	if err != nil {
//		util.Error("设置房间有效期失败[%v]", err.Error())
//		return err
//	}
//	return err
//}
//
////登记房间号
////格式:{"id":"1212", "owner":"wxUnionId", "time":new Date()}
//func InsertRoomNoData(roomId, ownerUnionId string) error {
//	err := MongoDB.C(CollectionRoomCache).Insert(bson.M{"id": roomId, "owner": ownerUnionId, "time": time.Now().Format("2006-01-02 15:04:05")})
//	if err != nil {
//		util.Error("链接房间唯一索引失败[%v]", err.Error())
//		return err
//	}
//	return err
//}
//
//
//
////创建用户列表并设置用户进入房间过期
////格式:{"id":"wxUnionId", "roomId":"123123", "time":new Date()}
//func CreateUserCache() error {
//	//创建房间号唯一索引
//	err := MongoDB.C(CollectionUserCache).EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true})
//	if err != nil {
//		util.Error("用户唯一索引失败[%v]", err.Error())
//		return err
//	}
//	//设置过期时间
//	err = MongoDB.C("userIds").EnsureIndex(mgo.Index{Key: []string{"time"},
//		ExpireAfter: 24 * time.Hour})
//	if err != nil {
//		util.Error("设置用户失效失败[%v]", err.Error())
//		return err
//	}
//	return err
//}
//
////登记用户房间号
////格式:{"id":"wxUnionId", "roomId":"123123", "time":time.now()}
//func InsertUserIdData(unionId, roomId string) error {
//	err := MongoDB.C(CollectionUserCache).Insert(bson.M{"id": unionId, "roomId": roomId, "time": time.Now().Format("2006-01-02 15:04:05")})
//	if err != nil {
//		util.Error("链接房间唯一索引失败[%v]", err.Error())
//	}
//	return err
//}
//
////查询用户房间号
//func QueryUserRoomId(unionId string)(roomId string, err error)  {
//
//	return roomId, err
//}
