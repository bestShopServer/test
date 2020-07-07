package zeromq

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/util"
	"encoding/json"
	"github.com/pebbe/zmq4"
	"sync"
	"time"
)

func InitZeroMQOneClient() (sock *zmq4.Socket, err error) {
	sock, err = zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return sock, err
	}
	if sock == nil {
		util.Error("MQSock ERROR: is null")
		return sock, err
	}
	err = sock.Connect(config.GetConfig().DBAddr)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	err = sock.SetConnectTimeout(time.Duration(config.GetConfig().DBTimeout))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
	}
	util.Debug("init socket connect ...")

	return sock, err
}

func TaskComm(sock *zmq4.Socket, data []byte) (msg []byte, err error) {
	util.Debug("TaskComm...")
	var mutex sync.Mutex

	mutex.Lock()
	l, err := sock.SendBytes(data, 0)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return nil, err
	}
	util.Debug("SendFrame len %v ...", l)

	msg, err = sock.RecvBytes(0)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return nil, err
	}
	util.Debug("RecvFrame ...")
	mutex.Unlock()

	return msg, err
}

func TaskJsonComm(sock *zmq4.Socket, data interface{}) (jsonResp map[string]interface{}, err error) {
	util.Debug("data[%v]", data)

	jraw, err := json.Marshal(data)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return jsonResp, err
	}
	res, err := TaskComm(sock, jraw)
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
