package zeromq

import (
	"DetectiveMasterServer/util"
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

type MQComm struct {
	Sock  *zmq.Socket
	Mutex sync.Mutex
	Wait  bool
}

var SockPool []MQComm

func InitMQSockPool(addr string, num int) {
	util.Info("InitMQSockPool...")
	var err error
	for i := 0; i < num; i++ {
		var sock MQComm
		sock.Sock, err = zmq.NewSocket(zmq.REQ)
		zmq.GetMaxSockets()
		if err != nil {
			util.Error("ERROR:%v", err.Error())
		}
		if sock.Sock != nil {
			err = sock.Sock.Connect(addr)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
			}
			sock.Wait = false
			util.Debug("init socket connect ...")
			SockPool = append(SockPool, sock)
		} else {
			util.Error("MQSock ERROR: is null")
			return
		}
	}
}

func ModifyMQSockPool(s *[]MQComm, index int, value MQComm) {
	var mutex sync.Mutex
	mutex.Lock()
	rear := append([]MQComm{}, (*s)[index+1:]...)
	*s = append(append((*s)[:index], value), rear...)
	mutex.Unlock()
}

func PoolGet() (sock MQComm, idx int, err error) {
	err = fmt.Errorf("未发现可用客户端")
	for idx, sock = range SockPool {
		if sock.Wait == false {
			sock.Wait = true
			return sock, idx, nil
		}
	}
	sock, idx, err = PoolGet()
	return sock, idx, err
}

func SendRecvDBMessage(data []byte) (reply []byte, err error) {
	util.Info("SendRecvDBMessage...")
	sock, idx, err := PoolGet()
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return reply, err
	}

	sock.Mutex.Lock()
	l, err := sock.Sock.SendBytes(data, 0)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return reply, err
	}
	util.Info("发送数据长度:%v", l)

	reply, err = sock.Sock.RecvBytes(0)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return reply, err
	}
	util.Debug("RecvBytes...")
	sock.Mutex.Unlock()
	sock.Wait = false
	ModifyMQSockPool(&SockPool, idx, sock)

	return reply, err
}

func TaskJsonNew(data interface{}) (jsonResp map[string]interface{}, err error) {
	util.Debug("data[%v]", data)

	jraw, err := json.Marshal(data)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return jsonResp, err
	}

	res, err := SendRecvDBMessage(jraw)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return jsonResp, err
	}

	err = json.Unmarshal(res, &jsonResp)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		return jsonResp, err
	}
	util.Debug("jsonResp...")

	return jsonResp, err
}
